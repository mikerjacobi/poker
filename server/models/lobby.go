package models

import (
	"errors"
	"fmt"
	"strings"

	"github.com/Sirupsen/logrus"
	"github.com/pborman/uuid"
	"github.com/spf13/viper"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

var (
	//lobby actions
	GameCreate   = "GAMECREATE"
	GameStart    = "GAMESTART"
	GameJoin     = "GAMEJOIN"
	GameLeave    = "GAMELEAVE"
	LobbyActions = []string{GameCreate, GameStart, GameJoin, GameLeave}
	GameTypes    []string

	//errors
	PlayerAlreadyJoined = errors.New("player already joined")
)

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

type LobbyMessage struct {
	Message
	Game `json:"game"`
}

type LobbyQueue struct {
	DB *mgo.Database
	Q  chan LobbyMessage
	*Comms
}

func NewLobbyQueue(db *mgo.Database, comms *Comms) (LobbyQueue, error) {
	GameTypes = strings.Split(viper.GetString("game_types"), ",")
	lq := LobbyQueue{
		DB:    db,
		Comms: comms,
	}

	lq.Q = make(chan LobbyMessage)
	go lq.ReadMessages()
	return lq, nil
}

func (lq LobbyQueue) ReadMessages() {
	for {
		lobbyMessage := <-lq.Q
		switch lobbyMessage.Type {
		case GameCreate:
			g, err := CreateGame(lq.DB, lobbyMessage.Game.Name, lobbyMessage.Game.GameType)
			if err != nil {
				logrus.Errorf("failed to create game: %s", err)
				continue
			}
			lq.SendAll(LobbyMessage{
				Message: lobbyMessage.Message,
				Game:    g,
			})
		case GameStart:
			logrus.Infof("game start in lobbyQ readmsgs")
		case GameJoin:
			game, err := JoinGame(lq.DB, lobbyMessage.Game.ID, lobbyMessage.Message.Sender)
			if err != nil {
				logrus.Errorf("failed to join game: %s", err)
				continue
			}
			lq.SendAll(LobbyMessage{
				Message: lobbyMessage.Message,
				Game:    game,
			})
		case GameLeave:
			accountID := lobbyMessage.Message.Sender.AccountID
			game, err := LeaveGame(lq.DB, lobbyMessage.Game.ID, accountID)
			if err != nil {
				logrus.Errorf("failed to leave game: %s", err)
				continue
			}
			lq.SendAll(LobbyMessage{
				Message: lobbyMessage.Message,
				Game:    game,
			})
		default:
			continue
		}
	}
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

	if name == "" {
		return Game{}, errors.New("gamename cannot be empty")
	}

	if !StringInSlice(gameType, GameTypes) {
		return Game{}, fmt.Errorf("invalid gametype: %s", gameType)
	}

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
			return Game{}, PlayerAlreadyJoined
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
			break
		}
	}
	if err := games.Update(query, g); err != nil {
		return Game{}, err
	}
	return g, nil
}
