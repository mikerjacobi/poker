package models

import (
	"errors"
	"fmt"

	"github.com/pborman/uuid"
	"github.com/spf13/viper"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

type GamePlayer struct {
	AccountID string `json:"accountID" bson:"accountID"`
	Name      string `json:"name" bson:"name"`
	Chips     int    `json:"chips" bson:"chips"`
}

type Game struct {
	ID   string `json:"gameID" bson:"gameID"`
	Name string `json:"gameName" bson:"gameName"`
	//State    string       `json:"state" bson:"state"`
	Players  []GamePlayer `json:"players" bson:"players"`
	GameType string       `json:"gameType" bson:"gameType"`
}

func CreateGame(db *mgo.Database, name, gameType string) (Game, error) {
	games := db.C("games")

	openGames, err := LoadGames(db)
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
	g := Game{gameID, name, []GamePlayer{}, gameType}
	if err := games.Insert(g); err != nil {
		return Game{}, fmt.Errorf("failed to insert: %s", err)
	}

	return g, nil
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
func LoadGames(db *mgo.Database) ([]Game, error) {
	gamesdb := db.C("games")
	games := []Game{}
	if err := gamesdb.Find(bson.M{}).All(&games); err != nil {
		return []Game{}, err
	}
	return games, nil
}

func (g Game) Update(db *mgo.Database) error {
	query := bson.M{"gameID": g.ID}
	if err := db.C("games").Update(query, g); err != nil {
		return err
	}
	return nil
}
