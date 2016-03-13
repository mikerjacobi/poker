package controllers

import (
	"encoding/json"
	"fmt"

	"github.com/mikerjacobi/poker/server/models"
	"gopkg.in/mgo.v2"
)

func HandleStart(msg models.Message) error {
	game := models.Game{}
	db := msg.Context.Get("db").(*mgo.Database)
	if err := json.Unmarshal(msg.Raw, &game); err != nil {
		return fmt.Errorf("failed to unmarshal game in highcardstart")
	}

	hcg, err := models.ToHighCardGame(game)
	if err != nil {
		return fmt.Errorf("failed to init high card game: %+v", err)
	}
	models.CreateHighCardGame(hcg)

	if err := hcg.PlayHand(db); err != nil {
		return fmt.Errorf("failed to start high card game: %+v", err)
	}
	return nil
}

func HandleReplay(msg models.Message) error {
	db := msg.Context.Get("db").(*mgo.Database)
	m := models.HighCardMessage{}
	if err := json.Unmarshal(msg.Raw, &m); err != nil {
		return fmt.Errorf("failed to unmarshal game in handle replay")
	}

	hcg, ok := models.GetHighCardGame(m.Game.ID)
	if !ok {
		return fmt.Errorf("failed to load high card game.  either not started or ended")
	}

	if err := hcg.PlayHand(db); err != nil {
		return fmt.Errorf("failed to start high card game: %+v", err)
	}
	return nil
}
