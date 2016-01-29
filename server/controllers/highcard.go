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

func (hcc HighCardController) ReadMessages() {
	for {
		m := <-hcc.Queue
		switch m.Type {
		case models.HighCardStart:
			if err := hcc.HandleStart(m.Raw); err != nil {
				logrus.Errorf("failed to start high card game: %+v", err)
				continue
			}
		case models.HighCardReplay:
			if err := hcc.HandleReplay(m.Raw); err != nil {
				logrus.Errorf("failed to replay high card game: %+v", err)
				continue
			}
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

func (hcc HighCardController) HandleStart(msg []byte) error {
	game := models.Game{}
	if err := json.Unmarshal(msg, &game); err != nil {
		return fmt.Errorf("failed to unmarshal game in highcardstart")
	}
	hcg, err := models.NewHighCardGame(game, hcc.Comms)
	if err != nil {
		return fmt.Errorf("failed to init high card game: %+v", err)
	}
	hcc.Games[game.ID] = hcg
	if err = hcc.Games[game.ID].Start(); err != nil {
		return fmt.Errorf("failed to start high card game: %+v", err)
	}
	return nil
}
func (hcc HighCardController) HandleReplay(msg []byte) error {
	game := models.HighCardMessage{}
	if err := json.Unmarshal(msg, &game); err != nil {
		return fmt.Errorf("failed to unmarshal game in handle replay")
	}

	hcg, ok := hcc.Games[game.Game.ID]
	if !ok {
		return fmt.Errorf("failed to load high card game.")
	}

	if err := hcg.PlayHand(); err != nil {
		return fmt.Errorf("failed to start high card game: %+v", err)
	}
	return nil
}
