package models

import (
	"errors"

	"github.com/pborman/uuid"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

type Game struct {
	GameID  string   `json:"game_id" bson:"game_id"`
	Name    string   `json:"name" bson:"name"`
	State   string   `json:"state" bson:"state"`
	Players []string `json:"players" bson:"players"`
}

var (
	PlayerAlreadyJoined = errors.New("player already joined")
)

func LoadGame(db *mgo.Database, gameID, gameName string) (Game, error) {
	gamesdb := db.C("games")
	game := Game{}
	var query bson.M
	if gameID != "" {
		query = bson.M{"game_id": gameID}
	} else if gameName != "" {
		query = bson.M{"name": gameName}
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
	gameID := uuid.New()
	g := Game{gameID, name, "open", []string{}}
	if err := games.Insert(g); err != nil {
		return Game{}, err
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
