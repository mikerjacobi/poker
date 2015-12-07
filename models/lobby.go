package models

import (
	"errors"
	"fmt"

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

	//errors
	PlayerAlreadyJoined = errors.New("player already joined")
)

type Game struct {
	ID      string   `json:"game_id" bson:"game_id"`
	Name    string   `json:"game_name" bson:"game_name"`
	State   string   `json:"state" bson:"state"`
	Players []string `json:"players" bson:"players"`
}

type LobbyMessage struct {
	Message
	Game `json:"game"`
}

type LobbyQueue struct {
	DB *mgo.Database
	Q  chan LobbyMessage
	CQ *CommsQueue
}

func NewLobbyQueue(db *mgo.Database, cq *CommsQueue) (LobbyQueue, error) {
	lq := LobbyQueue{
		DB: db,
		CQ: cq,
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
			g, err := CreateGame(lq.DB, lobbyMessage.Game.Name)
			if err != nil {
				logrus.Errorf("failed to create game: %s", err)
				continue
			}
			lq.CQ.SendAll(LobbyMessage{
				Message: lobbyMessage.Message,
				Game:    g,
			})
		case GameStart:
			logrus.Infof("game start in lobbyQ readmsgs")
		case GameJoin:
			accountID := lobbyMessage.Message.Sender.AccountID
			game, err := JoinGame(lq.DB, lobbyMessage.Game.ID, accountID)
			if err != nil {
				logrus.Errorf("failed to join game: %s", err)
				continue
			}
			lq.CQ.SendAll(LobbyMessage{
				Message: lobbyMessage.Message,
				Game:    game,
			})
		case GameLeave:
			logrus.Infof("game leave in lobbyQ readmsgs")
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
		query = bson.M{"game_id": gameID}
	} else if gameName != "" {
		query = bson.M{"game_name": gameName}
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
func CreateGame(db *mgo.Database, name string) (Game, error) {
	games := db.C("games")

	if name == "" {
		return Game{}, errors.New("gamename cannot be empty")
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
	g := Game{gameID, name, "open", []string{}}
	if err := games.Insert(g); err != nil {
		return Game{}, fmt.Errorf("failed to insert: %s", err)
	}
	return g, nil
}

func JoinGame(db *mgo.Database, gameID string, accountID string) (Game, error) {
	games := db.C("games")
	g := Game{}
	query := bson.M{"game_id": gameID}
	if err := games.Find(query).One(&g); err != nil {
		return Game{}, err
	}
	for i := range g.Players {
		if g.Players[i] == accountID {
			return Game{}, PlayerAlreadyJoined
		}
	}
	g.Players = append(g.Players, accountID)
	if err := games.Update(query, g); err != nil {
		return Game{}, err
	}
	return g, nil
}
