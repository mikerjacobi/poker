package controllers

import (
	"encoding/json"
	"fmt"
	"github.com/mikerjacobi/poker/server/models"
)

func CheckStartHoldem(game models.Game) error {
	if len(game.Players) < 2 {
		return fmt.Errorf("holdem: too few players")
	}
	gameJSON, err := json.Marshal(game)
	if err != nil {
		return fmt.Errorf("holdem: jsonmarshal error")
	}
	_ = models.Message{Type: models.HoldemStart, Raw: gameJSON}
	return nil
}
