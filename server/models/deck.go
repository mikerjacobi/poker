package models

import (
	"errors"
	"fmt"
	"github.com/Sirupsen/logrus"
	"math/rand"
	"time"
)

type Suit string

var hearts Suit = "hearts"
var diamonds Suit = "diamonds"
var clubs Suit = "clubs"
var spades Suit = "spades"
var suits = []Suit{hearts, diamonds, clubs, spades}

type Rank struct {
	Value   int
	Display string
}

var ace = Rank{14, "ace"}
var king = Rank{13, "king"}
var queen = Rank{12, "queen"}
var jack = Rank{11, "jack"}
var ten = Rank{10, "ten"}
var nine = Rank{9, "nine"}
var eight = Rank{8, "eight"}
var seven = Rank{7, "seven"}
var six = Rank{6, "six"}
var five = Rank{5, "five"}
var four = Rank{4, "four"}
var three = Rank{3, "three"}
var two = Rank{2, "two"}
var ranks = []Rank{ace, king, queen, jack, ten, nine, eight, seven, six, five, four, three, two}

type Card struct {
	Suit    `json:"suit"`
	Rank    int `json:"rank"`
	Display string
}

type Deck struct {
	Cards []Card
	Index int
}

func NewDeck() *Deck {
	logrus.Infof("new deck")
	d := Deck{
		Cards: make([]Card, 52),
		Index: 0,
	}
	k := 0
	for _, rank := range ranks {
		for _, suit := range suits {
			c := Card{
				Rank:    rank.Value,
				Suit:    suit,
				Display: fmt.Sprintf("%s of %s", rank.Display, suit),
			}
			d.Cards[k] = c
			k += 1
		}
	}
	d.Shuffle(3)
	return &d
}

func (d *Deck) Shuffle(shuffles int) {
	for i := 0; i < shuffles; i++ {
		for j := range d.Cards {
			r := rand.New(rand.NewSource(time.Now().UnixNano()))
			randInt := r.Intn(52)
			tmp := d.Cards[randInt]
			d.Cards[randInt] = d.Cards[j]
			d.Cards[j] = tmp
		}
	}
}

func (d *Deck) Deal() (*Card, error) {
	if d.Index == 52 {
		return nil, errors.New("no more cards in deck")
	}
	c := d.Cards[d.Index]
	d.Index += 1
	return &c, nil
}
