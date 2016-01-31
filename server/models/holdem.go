package models

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

type Holdem struct {
	ID      string   `json:"gameID" bson:"gameID"`
	Bank int `json:"-" bson:"bank"`
	Hands []Hand `json:"players" bson:"players"`
}
*/
var (
	//holdem actions
	HoldemStart   = "HOLDEMSTART"
	HoldemActions = []string{HoldemStart}
)

type HoldemMessage struct {
	Message
	Game `json:"game"`
}
