package models

import (
	"errors"
	"fmt"

	"github.com/pborman/uuid"
	"github.com/spf13/viper"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

var (
	//lobby actions
	GameCreate = "GAMECREATE"
	GameStart  = "GAMESTART"
	GameJoin   = "GAMEJOIN"

	GameJoinAlert  = "GAMEJOINALERT"
	GameLeave      = "GAMELEAVE"
	GameLeaveAlert = "GAMELEAVEALERT"
	LobbyActions   = []string{GameCreate, GameStart, GameJoin, GameJoinAlert, GameLeave, GameLeaveAlert}
)

type LobbyMessage struct {
	Message
	Game `json:"game"`
}
type GamePlayer struct {
	AccountID string `json:"accountID" bson:"accountID"`
	Name      string `json:"name" bson:"name"`
}

type Game struct {
	ID       string       `json:"gameID" bson:"gameID"`
	Name     string       `json:"gameName" bson:"gameName"`
	State    string       `json:"state" bson:"state"`
	Players  []GamePlayer `json:"players" bson:"players"`
	GameType string       `json:"gameType" bson:"gameType"`
}

func PlayerAccountIDs(gps []GamePlayer) []string {
	accountIDs := make([]string, len(gps))
	for i, gp := range gps {
		accountIDs[i] = gp.AccountID
	}
	return accountIDs
}

func LoadGame(db *mgo.Database, gameID, gameName string) (Game, error) {
	gamesdb := db.C("games")
	game := Game{}
	var query bson.M
	if gameID != "" {
		query = bson.M{"gameID": gameID}
	} else if gameName != "" {
		query = bson.M{"gameName": gameName}
	} else {
		return game, errors.New("gameid or gamename must be provided")
	}

	if err := gamesdb.Find(query).One(&game); err != nil {
		return Game{}, err
	}
	return game, nil
}
func LoadOpenGames(db *mgo.Database) ([]Game, error) {
	gamesdb := db.C("games")
	games := []Game{}
	query := bson.M{"state": "open"}
	if err := gamesdb.Find(query).All(&games); err != nil {
		return []Game{}, err
	}
	return games, nil
}
func CreateGame(db *mgo.Database, name, gameType string) (Game, error) {
	games := db.C("games")

	openGames, err := LoadOpenGames(db)
	maxOpenGames := viper.GetInt("max_open_games")
	if err == nil && len(openGames) >= maxOpenGames {
		return Game{}, fmt.Errorf("max_open_games limit reached: %d", maxOpenGames)
	} else if err != nil {
		return Game{}, fmt.Errorf("db error loading open games in create game: %s", err.Error())
	}

	_, err = LoadGame(db, "", name)
	if err == nil {
		return Game{}, fmt.Errorf("game name taken: %s", name)
	} else if err != mgo.ErrNotFound && err != nil {
		return Game{}, fmt.Errorf("db error loading game in create game: %s", err.Error())
	}

	gameID := uuid.New()
	g := Game{gameID, name, "open", []GamePlayer{}, gameType}
	if err := games.Insert(g); err != nil {
		return Game{}, fmt.Errorf("failed to insert: %s", err)
	}
	return g, nil
}

func JoinGame(db *mgo.Database, gameID string, account Account) (Game, error) {
	games := db.C("games")
	g := Game{}
	query := bson.M{"gameID": gameID}
	if err := games.Find(query).One(&g); err != nil {
		return Game{}, err
	}
	for i := range g.Players {
		if g.Players[i].AccountID == account.AccountID {
			return g, nil
		}
	}

	gp := GamePlayer{account.AccountID, account.Username}
	g.Players = append(g.Players, gp)
	if err := games.Update(query, g); err != nil {
		return Game{}, err
	}
	return g, nil
}

func LeaveGame(db *mgo.Database, gameID string, accountID string) (Game, error) {
	games := db.C("games")
	g := Game{}
	query := bson.M{"gameID": gameID}
	if err := games.Find(query).One(&g); err != nil {
		return Game{}, err
	}
	for i := range g.Players {
		if g.Players[i].AccountID == accountID {
			g.Players = append(g.Players[0:i], g.Players[i+1:]...)
			if err := games.Update(query, g); err != nil {
				return Game{}, err
			}
			break
		}
	}
	return g, nil
}

func RemovePlayerFromGames(db *mgo.Database, accountID string) error {
	gamesDB := db.C("games")
	games := []Game{}
	query := bson.M{"state": "open"}
	if err := gamesDB.Find(query).All(&games); err != nil {
		return err
	}

	for _, game := range games {
		for i := range game.Players {
			if game.Players[i].AccountID == accountID {
				game.Players = append(game.Players[0:i], game.Players[i+1:]...)
				if err := gamesDB.Update(bson.M{"gameID": game.ID}, game); err != nil {
					return err
				}
				break
			}
		}

	}
	return nil
}
