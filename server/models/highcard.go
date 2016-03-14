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
	Players  []HighCardPlayer `json:"players"`
	*Deck    `json:"-"`
	Pot      int `json:"pot"`
	ActionTo `json:"actionTo"`
}

type HighCardPlayer struct {
	GamePlayer `json:"gamePlayer"`
	Card       `json:"card"`
}

type HighCardGame struct {
	Game
	Hands []*Hand
	*Hand
	Ante int
}

type ActionTo struct {
	CallAmount int    `json:"callAmount"`
	AccountID  string `json:"accountID"`
}

func ToHighCardGame(game Game) *HighCardGame {
	hcg := HighCardGame{
		Game:  game,
		Hands: []*Hand{},
		Ante:  1,
	}
	return &hcg
}

func (hcg *HighCardGame) NewHand(db *mgo.Database) error {
	//update hcg with the current list of players
	game, err := LoadGame(db, hcg.Game.ID, "")
	if err != nil {
		return fmt.Errorf("failed to load game in newhand: %+v", err)
	}
	hcg.Game = game
	if !hcg.HighCardPlayable(db) {
		return fmt.Errorf("highcard game not in playable state")
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
		return fmt.Errorf("not enough players with ante")
	}
	if err := hcg.Game.Update(db); err != nil {
		return fmt.Errorf("failed to update game in new hand")
	}

	h.ActionTo = ActionTo{AccountID: h.Players[0].AccountID}

	hcg.Hand = &h
	hcg.Hands = append(hcg.Hands, &h)
	return nil
}

func (hcg *HighCardGame) StartHand(db *mgo.Database) error {
	//deal
	if err := hcg.NewHand(db); err != nil {
		return fmt.Errorf("failed to create new highcard hand: %+v", err)
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

func (hcg *HighCardGame) Send() error {
	for _, player := range hcg.Game.Players {
		//obscure all other players's cards
		obscuredPlayers := make([]HighCardPlayer, len(hcg.Hand.Players))
		copy(obscuredPlayers, hcg.Hand.Players)
		for i := range hcg.Hand.Players {
			//check and skip current player
			if hcg.Hand.Players[i].GamePlayer.AccountID == player.AccountID {
				continue
			}
			obscuredPlayers[i].Card = nullCard
		}

		state := Hand{
			Players:  obscuredPlayers,
			Pot:      hcg.Pot,
			ActionTo: hcg.ActionTo,
		}
		msg := HighCardMessage{
			Type:  highCardUpdate,
			Game:  hcg.Game,
			State: state,
		}
		if err := Send(player.AccountID, msg); err != nil {
			return fmt.Errorf("failed to send in highcard.start: %+v", err)
		}

	}
	return nil
}
