package controllers

import (
	"encoding/json"
	"fmt"
	"github.com/mikerjacobi/poker/server/models"
	"gopkg.in/mgo.v2"
)

func CheckStartHighCard(game models.Game) error {
	hcg, err := models.ToHighCardGame(game)
	if err != nil {
		return fmt.Errorf("failed to init high card game: %+v", err)
	}

	if !hcg.HighCardPlayable() {
		return fmt.Errorf("highcard game not in playable state")
	}
	gameJSON, err := json.Marshal(game)
	if err != nil {
		return fmt.Errorf("highcard: jsonmarshal error")
	}
	_ = models.Message{Type: models.HighCardStart, Raw: gameJSON}
	return nil
}

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
	game := models.HighCardMessage{}
	if err := json.Unmarshal(msg.Raw, &game); err != nil {
		return fmt.Errorf("failed to unmarshal game in handle replay")
	}

	hcg, ok := models.GetHighCardGame(game.Game.ID)
	if !ok {
		return fmt.Errorf("failed to load high card game.  either not started or ended")
	}

	if err := hcg.PlayHand(db); err != nil {
		return fmt.Errorf("failed to start high card game: %+v", err)
	}
	return nil
}
