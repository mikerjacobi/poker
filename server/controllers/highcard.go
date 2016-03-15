package controllers

import (
	"encoding/json"
	"fmt"

	"github.com/mikerjacobi/poker/server/models"
	"gopkg.in/mgo.v2"
)

type BetRaiseMessage struct {
	Amount int `json:"amount"`
	GameMessage
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

	if err := hcg.StartHand(db); err != nil {
		return fmt.Errorf("failed to start high card game: %+v", err)
	}

	if err := hcg.Send(); err != nil {
		return fmt.Errorf("failed to send highcard game: %+v", err)
	}

	return nil
}

func HandleCheck(msg models.Message) error {
	db := msg.Context.Get("db").(*mgo.Database)
	gm := GameMessage{}
	if err := json.Unmarshal(msg.Raw, &gm); err != nil {
		return fmt.Errorf("failed to unmarshal game in handle play")
	}

	account, ok := msg.Context.Get("user").(models.Account)
	if !ok {
		return fmt.Errorf("failed to get user in highcard check")
	}

	hcg, ok := models.GetHighCardGame(db, gm.GameID)
	if !ok {
		return fmt.Errorf("failed to load high card game.  either not started or ended")
	}

	if hcg.ActionTo.AccountID != account.AccountID {
		return fmt.Errorf("user attempted to act out of turn")
	}

	if err := hcg.Check(); err != nil {
		return fmt.Errorf("failed to start high card game: %+v", err)
	}

	if err := hcg.Send(); err != nil {
		return fmt.Errorf("failed to send highcard game: %+v", err)
	}

	return nil
}

func HandleBet(msg models.Message) error {
	db := msg.Context.Get("db").(*mgo.Database)
	gm := BetRaiseMessage{}
	if err := json.Unmarshal(msg.Raw, &gm); err != nil {
		return fmt.Errorf("failed to unmarshal game in handle play")
	}

	account, ok := msg.Context.Get("user").(models.Account)
	if !ok {
		return fmt.Errorf("failed to get user in highcard check")
	}

	hcg, ok := models.GetHighCardGame(db, gm.GameID)
	if !ok {
		return fmt.Errorf("failed to load high card game.  either not started or ended")
	}

	if hcg.ActionTo.AccountID != account.AccountID {
		return fmt.Errorf("user attempted to act out of turn")
	}

	if gm.Amount <= 0 {
		return fmt.Errorf("invalid bet amount passed in")
	}

	if hcg.Hand.Players[account.AccountID].Chips-gm.Amount < 0 {
		return fmt.Errorf("user attempted to bet more chips than owned")
	}

	if err := hcg.Bet(gm.Amount); err != nil {
		return fmt.Errorf("failed to start high card game: %+v", err)
	}

	if err := hcg.Send(); err != nil {
		return fmt.Errorf("failed to send highcard game: %+v", err)
	}

	return nil
}
