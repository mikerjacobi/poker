package models

import (
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
	//highcard actions
	HighCardStart   = "HIGHCARDSTART"
	HighCardActions = []string{HighCardStart}
)

type HighCardMessage struct {
	Message
	Game `json:"game"`
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
	logrus.Infof("game start hcg start")
	//deal
	h := hcg.NewHand()
	c, _ := h.Deck.Deal()
	logrus.Infof("card: %+v", c.Display)
	//hcg.SendPlayers(c)
	//msg := LobbyMessage{Message: msg, Game: game}
	//if err := hcg.Comms.SendGroup(msg, PlayerAccountIDs(hcg.Hand.Players)); err != nil{
	//	return fmt.Errorf("failed to sendgroup in highcard.start: %+v", err)
	//}
	//wait for action
	//    receive action
	//    validate action
	//    analyze gamestate
	// 		if endstate: goto deal
	//    else: emit action update
	return nil
}
