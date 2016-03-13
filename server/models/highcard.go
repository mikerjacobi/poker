package models

import (
	"fmt"

	"gopkg.in/mgo.v2"
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

var highCardManager HighCardManager

type HighCardManager struct {
	Games map[string]*HighCardGame
	DB    *mgo.Database
}

func InitializeHighCardManager(db *mgo.Database) {
	highCardManager = HighCardManager{
		DB:    db,
		Games: map[string]*HighCardGame{},
	}
}

func CreateHighCardGame(game *HighCardGame) {
	highCardManager.Games[game.ID] = game
}
func GetHighCardGame(gameID string) (*HighCardGame, bool) {
	hcg, ok := highCardManager.Games[gameID]
	return hcg, ok
}

var (
	//highcard server->client actions
	HighCardUpdate = "HIGHCARDUPDATE"
	HighCardError  = "HIGHCARDERROR"
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
	Hands []*Hand
	*Hand
}

func ToHighCardGame(game Game) (*HighCardGame, error) {
	hcg := HighCardGame{
		Game:  game,
		Hands: []*Hand{},
	}
	return &hcg, nil
}

func (hcg *HighCardGame) NewHand(db *mgo.Database) (*Hand, error) {
	//update hcg with the current list of players
	game, err := LoadGame(db, hcg.Game.ID, "")
	if err != nil {
		return nil, fmt.Errorf("failed to load game in newhand: %+v", err)
	}
	hcg.Game = game
	if !hcg.HighCardPlayable(db) {
		return nil, fmt.Errorf("highcard game not in playable state")
	}

	h := Hand{
		Players: hcg.Game.Players,
		Deck:    NewDeck(),
	}
	hcg.Hand = &h
	hcg.Hands = append(hcg.Hands, &h)
	return &h, nil
}

func (hcg *HighCardGame) PlayHand(db *mgo.Database) error {
	//deal
	hand, err := hcg.NewHand(db)
	if err != nil {
		return fmt.Errorf("failed to create new highcard hand: %+v", err)
	}
	card, _ := hand.Deck.Deal()
	msg := HighCardMessage{
		Message: Message{Type: HighCardUpdate},
		Game:    hcg.Game,
		State: struct {
			Card `json:"card"`
		}{*card},
	}

	if err := SendGame(db, hcg.Game.ID, msg); err != nil {
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

func (hcg *HighCardGame) HighCardPlayable(db *mgo.Database) bool {
	if len(hcg.Game.Players) < 2 {
		SendGameError(db, hcg.Game.ID, "not enough players to start game.")
		return false
	}
	//for _, p := range hcg.Game.Players {
	//	if p.Balance <= 0 {
	//		hcg.SendError("you need to buy in.", []string{p.AccountID})
	//		return false
	//	}
	//}
	return true
}
