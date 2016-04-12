package models

import (
	"fmt"

	"github.com/Sirupsen/logrus"

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
}

type HighCardMessage struct {
	Type  string `json:"type"`
	Game  `json:"gameInfo"`
	State interface{} `json:"gameState"`
}

type HighCardHand struct {
	PlayerList      []*HighCardPlayer          `json:"players"`
	Players         map[string]*HighCardPlayer `json:"-"`
	*Deck           `json:"-"`
	Pot             int `json:"pot"`
	ActionTo        `json:"actionTo"`
	Complete        bool `json:"complete"`
	NumTurns        int  `json:"-"`
	NumStartPlayers int  `json:"-"`
}

type HighCardPlayer struct {
	Next       *HighCardPlayer `json:"-"`
	GamePlayer `json:"gamePlayer"`
	Card       `json:"card"`
	State      string `json:"state"`
}

type HighCardGame struct {
	Game
	Hands []*HighCardHand
	Hand  *HighCardHand
	Ante  int
}

type ActionTo struct {
	CallAmount int    `json:"callAmount"`
	AccountID  string `json:"accountID"`
}

func InitializeHighCardManager() {
	highCardManager = HighCardManager{
		Games: map[string]*HighCardGame{},
	}
}

func LoadCreateHighCardGame(db *mgo.Database, gameID string) (*HighCardGame, error) {
	game, err := LoadGame(db, gameID, "")
	if err != nil {
		return nil, fmt.Errorf("failed to load game in get high card game: %+v", err)
	}

	hcg, ok := highCardManager.Games[gameID]
	if ok {
		hcg.Game = game
	} else {
		hcg = &HighCardGame{
			Game:  game,
			Hands: []*HighCardHand{},
			Ante:  1,
		}
		highCardManager.Games[game.ID] = hcg
	}

	return hcg, nil
}

func (hcp *HighCardPlayer) Copy() *HighCardPlayer {
	newPlayer := *hcp
	return &newPlayer
}

func (hcg *HighCardGame) NewHand() error {
	//update hcg with the current list of players
	if len(hcg.Game.Players) < 2 {
		return fmt.Errorf("not enough players to start game.")
	}

	h := HighCardHand{
		PlayerList: []*HighCardPlayer{},
		Players:    map[string]*HighCardPlayer{},
		Deck:       NewDeck(),
		Pot:        0,
		NumTurns:   0,
		Complete:   false,
	}

	for i, p := range hcg.Game.Players {
		if p.Chips-hcg.Ante < 0 {
			continue
		}
		hcg.Game.Players[i].Chips -= hcg.Ante
		h.Pot += hcg.Ante
		card, _ := h.Deck.Deal()

		player := HighCardPlayer{
			GamePlayer: hcg.Game.Players[i],
			Card:       *card,
			State:      "",
		}

		h.PlayerList = append(h.PlayerList, &player)
		h.Players[hcg.Game.Players[i].AccountID] = &player
	}

	if len(h.PlayerList) < 2 {
		return fmt.Errorf("not enough players with ante")
	}

	//initiate each player's next player
	for i := range h.PlayerList {
		if i == len(h.PlayerList)-1 {
			h.PlayerList[i].Next = h.PlayerList[0]
			break
		}
		h.PlayerList[i].Next = h.PlayerList[i+1]
	}
	h.NumStartPlayers = len(h.PlayerList)
	h.ActionTo = ActionTo{AccountID: h.PlayerList[0].AccountID}

	hcg.Hand = &h
	hcg.Hands = append(hcg.Hands, &h)
	return nil
}

func (hcg *HighCardGame) StartHand() error {
	//deal
	if err := hcg.NewHand(); err != nil {
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

func (hcg *HighCardGame) Send() error {
	for _, player := range hcg.Game.Players {
		//obscure all other players's cards
		obscuredPlayers := make([]*HighCardPlayer, len(hcg.Hand.PlayerList))
		for i := range hcg.Hand.PlayerList {
			//check and skip current player
			obscuredPlayers[i] = hcg.Hand.PlayerList[i].Copy()
			if hcg.Hand.PlayerList[i].GamePlayer.AccountID != player.AccountID {
				obscuredPlayers[i].Card = nullCard
			}
		}

		state := HighCardHand{
			PlayerList: obscuredPlayers,
			Pot:        hcg.Hand.Pot,
			ActionTo:   hcg.Hand.ActionTo,
			Complete:   hcg.Hand.Complete,
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

func (hcg *HighCardGame) Fold() error {
	if hcg.Hand.Complete {
		return nil
	}
	return hcg.Transition("fold")
}

func (hcg *HighCardGame) Check() error {
	if hcg.Hand.Complete {
		return nil
	}
	return hcg.Transition("check")
}

func (hcg *HighCardGame) Call() error {
	if hcg.Hand.Complete {
		return nil
	}
	currPlayer := hcg.Hand.ActionTo.AccountID
	hcg.Hand.Players[currPlayer].Chips -= hcg.Hand.ActionTo.CallAmount
	hcg.Hand.Pot += hcg.Hand.ActionTo.CallAmount
	hcg.Hand.ActionTo.CallAmount = 0
	return hcg.Transition("call")
}

func (hcg *HighCardGame) Bet(amount int) error {
	if hcg.Hand.Complete {
		return nil
	}
	currPlayer := hcg.Hand.ActionTo.AccountID
	hcg.Hand.Players[currPlayer].Chips -= amount
	hcg.Hand.Pot += amount
	hcg.Hand.ActionTo.CallAmount = amount
	return hcg.Transition(fmt.Sprintf("bet %d", amount))
}

func (hcg *HighCardGame) Raise(amount int) error {
	if hcg.Hand.Complete {
		return nil
	}
	currPlayer := hcg.Hand.ActionTo.AccountID
	hcg.Hand.Players[currPlayer].Chips -= (hcg.Hand.ActionTo.CallAmount + amount)
	hcg.Hand.Pot += (hcg.Hand.ActionTo.CallAmount + amount)
	hcg.Hand.ActionTo.CallAmount = amount
	return hcg.Transition(fmt.Sprintf("call %d, re-raise %d", hcg.Hand.ActionTo.CallAmount, amount))
}

func (hcg *HighCardGame) Transition(state string) error {
	currPlayer := hcg.Hand.ActionTo.AccountID
	hcg.Hand.Players[currPlayer].State = state

	//todo while loop to not be fold statea
	hcg.Hand.ActionTo.AccountID = hcg.Hand.Players[currPlayer].Next.GamePlayer.AccountID
	for {
		if hcg.Hand.ActionTo.AccountID == "fold" {
			currPlayer := hcg.Hand.ActionTo.AccountID
			hcg.Hand.ActionTo.AccountID = hcg.Hand.Players[currPlayer].Next.GamePlayer.AccountID
		} else {
			break
		}
	}

	return hcg.checkComplete()
}

func (hcg *HighCardGame) checkComplete() error {
	hcg.Hand.NumTurns++
	if hcg.Hand.NumTurns >= hcg.Hand.NumStartPlayers && hcg.Hand.ActionTo.CallAmount == 0 {
		hcg.Hand.Complete = true
	} else if hcg.OnePlayerRemains() {
		hcg.Hand.Complete = true
	}

	if !hcg.Hand.Complete {
		return nil
	}

	highest := 0
	winners := []string{}
	for accountID, player := range hcg.Hand.Players {
		if player.State == "fold" {
			continue
		}

		if player.Card.Rank == highest {
			winners = append(winners, accountID)
		} else if player.Card.Rank > highest {
			highest = player.Card.Rank
			winners = []string{accountID}
		}
	}

	//calculate this hands chip updates
	payout := hcg.Hand.Pot / len(winners)
	for _, w := range winners {
		hcg.Hand.Players[w].Chips += payout
		hcg.Hand.Players[w].State = "winner"
	}

	//record the chip updates in our game object which will get saved
	for i, p := range hcg.Game.Players {
		hcg.Game.Players[i].Chips = hcg.Hand.Players[p.AccountID].Chips

		//if you're not a winner you're a loser
		if hcg.Hand.Players[p.AccountID].State != "winner" {
			hcg.Hand.Players[p.AccountID].State = "loser"
		}
	}

	logrus.Infof("highcard game complete.  winners: %+v, payment: %d", winners, payout)
	hcg.Hand.Pot = 0
	return nil
}

func (hcg *HighCardGame) OnePlayerRemains() bool {
	//determine if only one player is in a nonfold state
	numActives := 0
	for _, p := range hcg.Hand.Players {
		if p.State != "fold" {
			numActives++
		}
	}
	if numActives == 1 {
		return true
	}
	return false
}
