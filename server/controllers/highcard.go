package controllers

import (
	"encoding/json"
	"fmt"
	"github.com/Sirupsen/logrus"
	"github.com/mikerjacobi/poker/server/models"
	"gopkg.in/mgo.v2"
)

var (
	//highcard actions
	HighCardStart   = "HIGHCARDSTART"
	HighCardActions = []string{HighCardStart}
)

type HighCardMessage struct {
	Message
	models.Game `json:"game"`
}

type HighCardController struct {
	DB    *mgo.Database
	Queue chan Message
	*models.Comms
}

func newHighCardController(db *mgo.Database, c *models.Comms) (HighCardController, error) {
	hc := HighCardController{
		DB:    db,
		Comms: c,
	}

	hc.Queue = make(chan Message)
	go hc.ReadMessages()
	return hc, nil
}

func (hc HighCardController) ReadMessages() {
	for {
		m := <-hc.Queue
		switch m.Type {
		case HighCardStart:
			logrus.Infof("game start in highcardQ readmsgs")
		default:
			continue
		}
	}
}

func (hcc HighCardController) CheckStartGame(game models.Game) error {
	if len(game.Players) < 2 {
		return fmt.Errorf("highcard: too few players")
	}
	gameJSON, err := json.Marshal(game)
	if err != nil {
		return fmt.Errorf("highcard: jsonmarshal error")
	}
	m := Message{Type: HighCardStart, Raw: gameJSON}
	hcc.Queue <- m
	return nil
}
