package controllers

import (
	"encoding/json"
	"fmt"
	"github.com/Sirupsen/logrus"
	"github.com/mikerjacobi/poker/server/models"
	"gopkg.in/mgo.v2"
)

type HighCardController struct {
	DB    *mgo.Database
	Queue chan models.Message
	*models.Comms
	Games map[string]*models.HighCardGame
}

func newHighCardController(db *mgo.Database, c *models.Comms) (HighCardController, error) {
	hc := HighCardController{
		DB:    db,
		Comms: c,
		Games: map[string]*models.HighCardGame{},
	}

	hc.Queue = make(chan models.Message)
	go hc.ReadMessages()
	return hc, nil
}

func (hc HighCardController) ReadMessages() {
	for {
		m := <-hc.Queue
		switch m.Type {
		case models.HighCardStart:
			game := models.Game{}
			if err := json.Unmarshal(m.Raw, &game); err != nil {
				logrus.Errorf("failed to unmarshal game in highcardstart")
				return
			}
			hcg, err := models.NewHighCardGame(game, hc.Comms)
			if err != nil {
				logrus.Errorf("failed to start high card game: %+v", err)
				continue
			}
			hc.Games[game.ID] = hcg
			hc.Games[game.ID].Start()
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
	m := models.Message{Type: models.HighCardStart, Raw: gameJSON}
	hcc.Queue <- m
	return nil
}
