package models

import (
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

type Counter struct {
	Count int `json:"count" bson:"count"`
}

func LoadMathCount(db *mgo.Database) (Counter, error) {
	mathdb := db.C("math")
	counter := Counter{}
	query := bson.M{}
	if err := mathdb.Find(query).One(&counter); err != nil {
		return counter, err
	}
	return counter, nil
}
