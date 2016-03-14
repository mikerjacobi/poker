package controllers

import (
	"encoding/json"
	"fmt"

	"github.com/mikerjacobi/poker/server/models"
	"gopkg.in/mgo.v2"
)

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

	if err := hcg.StartHand(db); err != nil {
		return fmt.Errorf("failed to start high card game: %+v", err)
	}
	return nil
}
