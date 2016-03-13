package controllers

import (
	"encoding/json"
	"fmt"

	"github.com/mikerjacobi/poker/server/models"
	"gopkg.in/mgo.v2"
)

const (
	highCardUpdate = "/highcard/update"
)

type HighCardMessage struct {
	Type        string `json:"type"`
	models.Game `json:"gameInfo"`
	State       interface{} `json:"gameState"`
}

func HandlePlay(msg models.Message) error {
	db := msg.Context.Get("db").(*mgo.Database)
	gm := GameMessage{}
	if err := json.Unmarshal(msg.Raw, &gm); err != nil {
		return fmt.Errorf("failed to unmarshal game in handle play")
	}

	hcg, ok := models.GetHighCardGame(db, gm.GameID)
	if !ok {
		return fmt.Errorf("failed to load high card game.  either not started or ended")
	}

	card, err := hcg.PlayHand(db)
	if err != nil {
		return fmt.Errorf("failed to start high card game: %+v", err)
	}

	hcMsg := HighCardMessage{
		Type: highCardUpdate,
		Game: hcg.Game,
		State: struct {
			models.Card `json:"card"`
		}{*card},
	}
	if err := models.SendGame(db, hcg.Game.ID, hcMsg); err != nil {
		return fmt.Errorf("failed to sendgroup in highcard.start: %+v", err)
	}

	return nil
}
