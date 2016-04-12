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

	hcg, err := models.LoadCreateHighCardGame(db, gm.GameID)
	if err != nil {
		return fmt.Errorf("failed to load high card game.  either not started or ended")
	}

	if len(hcg.Hands) > 0 && !hcg.Hand.Complete {
		return fmt.Errorf("cannot start new hand, game in progress")
	}

	if err := hcg.StartHand(); err != nil {
		return fmt.Errorf("failed to start high card game: %+v", err)
	}

	if err := hcg.Game.Update(db); err != nil {
		return fmt.Errorf("failed to update game in play")
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

	hcg, err := models.LoadCreateHighCardGame(db, gm.GameID)
	if err != nil {
		return fmt.Errorf("failed to load high card game.  either not started or ended")
	}

	if hcg.Hand.ActionTo.AccountID != account.AccountID {
		return fmt.Errorf("user attempted to act out of turn")
	}

	if err := hcg.Check(); err != nil {
		return fmt.Errorf("failed to start high card game: %+v", err)
	}

	if err := hcg.Game.Update(db); err != nil {
		return fmt.Errorf("failed to update game in check")
	}

	if err := hcg.Send(); err != nil {
		return fmt.Errorf("failed to send highcard game: %+v", err)
	}

	return nil
}

func HandleFold(msg models.Message) error {
	db := msg.Context.Get("db").(*mgo.Database)
	gm := GameMessage{}
	if err := json.Unmarshal(msg.Raw, &gm); err != nil {
		return fmt.Errorf("failed to unmarshal game in handle play")
	}

	account, ok := msg.Context.Get("user").(models.Account)
	if !ok {
		return fmt.Errorf("failed to get user in highcard check")
	}

	hcg, err := models.LoadCreateHighCardGame(db, gm.GameID)
	if err != nil {
		return fmt.Errorf("failed to load high card game.  either not started or ended")
	}

	if hcg.Hand.ActionTo.AccountID != account.AccountID {
		return fmt.Errorf("user attempted to act out of turn")
	}

	if err := hcg.Fold(); err != nil {
		return fmt.Errorf("failed to start high card game: %+v", err)
	}

	if err := hcg.Game.Update(db); err != nil {
		return fmt.Errorf("failed to update game in fold")
	}

	if err := hcg.Send(); err != nil {
		return fmt.Errorf("failed to send highcard game: %+v", err)
	}

	return nil
}

func HandleCall(msg models.Message) error {
	db := msg.Context.Get("db").(*mgo.Database)
	gm := GameMessage{}
	if err := json.Unmarshal(msg.Raw, &gm); err != nil {
		return fmt.Errorf("failed to unmarshal game in handle play")
	}

	account, ok := msg.Context.Get("user").(models.Account)
	if !ok {
		return fmt.Errorf("failed to get user in highcard check")
	}

	hcg, err := models.LoadCreateHighCardGame(db, gm.GameID)
	if err != nil {
		return fmt.Errorf("failed to load high card game.  either not started or ended")
	}

	if hcg.Hand.ActionTo.AccountID != account.AccountID {
		return fmt.Errorf("user attempted to act out of turn")
	}

	if hcg.Hand.Players[account.AccountID].Chips-hcg.Hand.ActionTo.CallAmount < 0 {
		return fmt.Errorf("user attempted to call more chips than owned")
	}

	if err := hcg.Call(); err != nil {
		return fmt.Errorf("failed to call high card game: %+v", err)
	}

	if err := hcg.Game.Update(db); err != nil {
		return fmt.Errorf("failed to update game in call")
	}

	if err := hcg.Game.Update(db); err != nil {
		return fmt.Errorf("failed to update game in call")
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

	hcg, err := models.LoadCreateHighCardGame(db, gm.GameID)
	if err != nil {
		return fmt.Errorf("failed to load high card game.  either not started or ended")
	}

	if hcg.Hand.ActionTo.AccountID != account.AccountID {
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

	if err := hcg.Game.Update(db); err != nil {
		return fmt.Errorf("failed to update game in bet")
	}

	if err := hcg.Send(); err != nil {
		return fmt.Errorf("failed to send highcard game: %+v", err)
	}

	return nil
}

func HandleRaise(msg models.Message) error {
	db := msg.Context.Get("db").(*mgo.Database)
	gm := BetRaiseMessage{}
	if err := json.Unmarshal(msg.Raw, &gm); err != nil {
		return fmt.Errorf("failed to unmarshal game in handle play")
	}

	account, ok := msg.Context.Get("user").(models.Account)
	if !ok {
		return fmt.Errorf("failed to get user in highcard check")
	}

	hcg, err := models.LoadCreateHighCardGame(db, gm.GameID)
	if err != nil {
		return fmt.Errorf("failed to load high card game.  either not started or ended")
	}

	if hcg.Hand.ActionTo.AccountID != account.AccountID {
		return fmt.Errorf("user attempted to act out of turn")
	}

	if gm.Amount <= 0 {
		return fmt.Errorf("invalid raise amount passed in")
	}

	if hcg.Hand.Players[account.AccountID].Chips-(hcg.Hand.ActionTo.CallAmount+gm.Amount) < 0 {
		return fmt.Errorf("user attempted to raise more chips than owned")
	}

	if err := hcg.Raise(gm.Amount); err != nil {
		return fmt.Errorf("failed to raise high card game: %+v", err)
	}

	if err := hcg.Game.Update(db); err != nil {
		return fmt.Errorf("failed to update game in raise")
	}

	if err := hcg.Send(); err != nil {
		return fmt.Errorf("failed to send highcard game: %+v", err)
	}

	return nil
}
