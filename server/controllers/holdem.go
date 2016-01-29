package controllers

import (
	"encoding/json"
	"fmt"
	"github.com/Sirupsen/logrus"
	"github.com/mikerjacobi/poker/server/models"
	"gopkg.in/mgo.v2"
)

type HoldemController struct {
	DB    *mgo.Database
	Queue chan models.Message
	*models.Comms
}

func newHoldemController(db *mgo.Database, c *models.Comms) (HoldemController, error) {
	hc := HoldemController{
		DB:    db,
		Comms: c,
	}

	hc.Queue = make(chan models.Message)
	go hc.ReadMessages()
	return hc, nil
}

func (hc HoldemController) ReadMessages() {
	for {
		m := <-hc.Queue
		switch m.Type {
		case models.GameStart:
			logrus.Infof("game start in holdemQ readmsgs")
		default:
			continue
		}
	}
}
func (hc HoldemController) CheckStartGame(game models.Game) error {
	if len(game.Players) < 2 {
		return fmt.Errorf("holdem: too few players")
	}
	gameJSON, err := json.Marshal(game)
	if err != nil {
		return fmt.Errorf("holdem: jsonmarshal error")
	}
	m := models.Message{Type: models.HoldemStart, Raw: gameJSON}
	hc.Queue <- m
	return nil
}
