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
	//game actions
	GameCreate  = "GAMECREATE"
	GameStart   = "GAMESTART"
	GameJoin    = "GAMEJOIN"
	GameLeave   = "GAMELEAVE"
	GameActions = []string{GameCreate, GameStart, GameJoin, GameLeave}

	//errors
	PlayerAlreadyJoined = errors.New("player already joined")
)

type Game struct {
	GameID  string   `json:"game_id" bson:"game_id"`
	Name    string   `json:"game_name" bson:"game_name"`
	State   string   `json:"state" bson:"state"`
	Players []string `json:"players" bson:"players"`
}

type GameMessage struct {
	Message
	Game `json:"game"`
}

type GameQueue struct {
	DB *mgo.Database
	Q  chan GameMessage
	CQ *CommsQueue
}

func NewGameQueue(db *mgo.Database, cq *CommsQueue) (GameQueue, error) {
	gq := GameQueue{
		DB: db,
		CQ: cq,
	}

	gq.Q = make(chan GameMessage)
	go gq.ReadMessages()
	return gq, nil
}

func (gq GameQueue) ReadMessages() {
	for {
		gameMessage := <-gq.Q
		switch gameMessage.Type {
		case GameCreate:
			g, err := CreateGame(gq.DB, gameMessage.Game.Name)
			if err != nil {
				logrus.Errorf("failed to create game: %s", err)
				continue
			}
			gq.CQ.SendAll(GameMessage{
				Message: gameMessage.Message,
				Game:    g,
			})
		case GameStart:
			logrus.Infof("game start in gameQ readmsgs")
		case GameJoin:
			logrus.Infof("game join in gameQ readmsgs")
			/*
				//func JoinGame(db *mgo.Database, gameID string, accountID string) error {
				if err != nil{
					logrus.Errorf("failed to create game: %s", err)
					continue
				}
				gq.CQ.SendAll(GameMessage{
					Message: gameMessage.Message,
					Game: g,
				})
			*/
		case GameLeave:
			logrus.Infof("game leave in gameQ readmsgs")
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

func JoinGame(db *mgo.Database, gameID string, accountID string) error {
	games := db.C("games")
	g := Game{}
	query := bson.M{"game_id": gameID}
	if err := games.Find(query).One(&g); err != nil {
		return err
	}
	for i := range g.Players {
		if g.Players[i] == accountID {
			return PlayerAlreadyJoined
		}
	}
	g.Players = append(g.Players, accountID)
	if err := games.Update(query, g); err != nil {
		return err
	}
	return nil
}
