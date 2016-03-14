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

const (
	highCardUpdate = "/highcard/update"
)

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

type HighCardMessage struct {
	Type  string `json:"type"`
	Game  `json:"gameInfo"`
	State interface{} `json:"gameState"`
}

func CreateHighCardGame(game *HighCardGame) {
	highCardManager.Games[game.ID] = game
}
func GetHighCardGame(db *mgo.Database, gameID string) (*HighCardGame, bool) {
	hcg, ok := highCardManager.Games[gameID]
	if !ok {
		game, err := LoadGame(db, gameID, "")
		if err != nil {
			return nil, false
		}
		hcg = ToHighCardGame(game)
		CreateHighCardGame(hcg)
	}
	return hcg, true
}

type Hand struct {
	Players []HighCardPlayer
	*Deck
	Pot int
}

type HighCardPlayer struct {
	GamePlayer `json:"game_player"`
	Card       `json:"card"`
	IsButton   bool `json:"is_button"`
}

type HighCardGame struct {
	Game
	Hands []*Hand
	*Hand
	Ante int
}

func ToHighCardGame(game Game) *HighCardGame {
	hcg := HighCardGame{
		Game:  game,
		Hands: []*Hand{},
		Ante:  1,
	}
	return &hcg
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
		Players: []HighCardPlayer{},
		Deck:    NewDeck(),
		Pot:     0,
	}
	for i, p := range hcg.Game.Players {
		if p.Chips-hcg.Ante < 0 {
			continue
		}
		hcg.Game.Players[i].Chips -= hcg.Ante
		h.Pot += hcg.Ante
		card, _ := h.Deck.Deal()
		h.Players = append(h.Players, HighCardPlayer{
			GamePlayer: hcg.Game.Players[i],
			Card:       *card,
		})
	}

	if len(h.Players) < 2 {
		return nil, fmt.Errorf("not enough players with ante")
	}
	if err := hcg.Game.Update(db); err != nil {
		return nil, fmt.Errorf("failed to update game in new hand")
	}

	hcg.Hand = &h
	hcg.Hands = append(hcg.Hands, &h)
	return &h, nil
}

func (hcg *HighCardGame) StartHand(db *mgo.Database) error {
	//deal
	hand, err := hcg.NewHand(db)
	if err != nil {
		return fmt.Errorf("failed to create new highcard hand: %+v", err)
	}

	for _, _ = range hand.Players {
		hcMsg := HighCardMessage{
			Type: highCardUpdate,
			Game: hcg.Game,
			State: struct {
				Players []HighCardPlayer `json:"players"`
			}{hand.Players},
		}

		//TODO obscure all other players's cards

		if err := SendGame(db, hcg.Game.ID, hcMsg); err != nil {
			return fmt.Errorf("failed to sendgroup in highcard.start: %+v", err)
		}

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
