package models

import (
	"github.com/Sirupsen/logrus"
	"gopkg.in/mgo.v2"
)

var (
	//holdem actions
	HoldemStart   = "HOLDEMSTART"
	HoldemActions = []string{HoldemStart}
)

/*
type Player struct{
	AccountID string
	Chips int
}

type State struct{
	CurrentPlayers []Player
	ActivePlayer string
}

type Bet

type Hand struct{
	Pot int
	Bank int
	State State
	Actions []interface
}

type Holdem struct {
	ID      string   `json:"gameID" bson:"gameID"`
	Bank int `json:"-" bson:"bank"`
	Hands []Hand `json:"players" bson:"players"`
}
*/

type HoldemMessage struct {
	Message
	Game `json:"game"`
}

type HoldemQueue struct {
	DB *mgo.Database
	Q  chan HoldemMessage
	*Comms
}

func NewHoldemQueue(db *mgo.Database, c *Comms) (HoldemQueue, error) {
	hq := HoldemQueue{
		DB:    db,
		Comms: c,
	}

	hq.Q = make(chan HoldemMessage)
	go hq.ReadMessages()
	return hq, nil
}

func (hq HoldemQueue) ReadMessages() {
	for {
		holdemMessage := <-hq.Q
		switch holdemMessage.Type {
		case GameStart:
			logrus.Infof("game start in holdemQ readmsgs")
		default:
			continue
		}
	}
}
