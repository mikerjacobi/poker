package models

import (
	"fmt"
	"github.com/Sirupsen/logrus"
)

/*
type Player struct{
	AccountID string
	Chips int
}

type State struct{
	CurrentPlayers []Player
	ActivePlayer string
}

type Bet

type Hand struct{
	Pot int
	Bank int
	State State
	Actions []interface
}

type HighCard struct {
	ID      string   `json:"gameID" bson:"gameID"`
	Bank int `json:"-" bson:"bank"`
	Hands []Hand `json:"players" bson:"players"`
}
*/

var (
	//highcard client->server actions
	HighCardStart   = "HIGHCARDSTART"
	HighCardReplay  = "HIGHCARDREPLAY"
	HighCardActions = []string{HighCardStart, HighCardReplay}

	//highcard server->client actions
	HighCardUpdate = "HIGHCARDUPDATE"
)

type HighCardMessage struct {
	Message
	Game  `json:"gameInfo"`
	State interface{} `json:"gameState"`
}

type Hand struct {
	Players []GamePlayer
	*Deck
}

type HighCardGame struct {
	Game
	*Comms
	Hands []*Hand
	*Hand
}

func NewHighCardGame(game Game, comms *Comms) (*HighCardGame, error) {
	hcg := HighCardGame{
		Game:  game,
		Comms: comms,
		Hands: []*Hand{},
	}
	return &hcg, nil
}

func (hcg *HighCardGame) NewHand() *Hand {
	h := Hand{
		Players: hcg.Game.Players,
		Deck:    NewDeck(),
	}
	hcg.Hand = &h
	hcg.Hands = append(hcg.Hands, &h)
	return &h
}

func (hcg *HighCardGame) Start() error {
	//deal
	h := hcg.NewHand()
	c, _ := h.Deck.Deal()
	logrus.Infof("card: %+v", c.Display)
	msg := HighCardMessage{
		Message: Message{Type: HighCardUpdate},
		Game:    hcg.Game,
		State: struct {
			Card `json:"card"`
		}{*c},
	}
	accountIDs := PlayerAccountIDs(hcg.Hand.Players)
	if err := hcg.Comms.SendGroup(msg, accountIDs); err != nil {
		return fmt.Errorf("failed to sendgroup in highcard.start: %+v", err)
	}
	//wait for action
	//    receive action
	//    validate action
	//    analyze gamestate
	// 		if endstate: goto deal
	//    else: emit action update
	return nil
}

func (hcg *HighCardGame) PlayHand() error {
	//deal
	h := hcg.NewHand()
	c, _ := h.Deck.Deal()
	msg := HighCardMessage{
		Message: Message{Type: HighCardUpdate},
		Game:    hcg.Game,
		State: struct {
			Card `json:"card"`
		}{*c},
	}
	accountIDs := PlayerAccountIDs(hcg.Hand.Players)
	if err := hcg.Comms.SendGroup(msg, accountIDs); err != nil {
		return fmt.Errorf("failed to sendgroup in highcard.start: %+v", err)
	}
	//wait for action
	//    receive action
	//    validate action
	//    analyze gamestate
	// 		if endstate: goto deal
	//    else: emit action update
	return nil
}
