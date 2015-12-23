package models

import (
	"math"

	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

type Counter struct {
	Count int `json:"count" bson:"count"`
}

func Increment(db *mgo.Database) (*Counter, error) {
	c, err := LoadMathCount(db)
	if err != nil {
		return nil, err
	}
	c.Count = checkBounds(c.Count + 1)
	if err := saveMathCount(db, c); err != nil {
		return nil, err
	}
	return c, nil
}
func Decrement(db *mgo.Database) (*Counter, error) {
	c, err := LoadMathCount(db)
	if err != nil {
		return nil, err
	}
	c.Count = checkBounds(c.Count - 1)
	if err := saveMathCount(db, c); err != nil {
		return nil, err
	}
	return c, nil
}
func Square(db *mgo.Database) (*Counter, error) {
	c, err := LoadMathCount(db)
	if err != nil {
		return nil, err
	}
	c.Count = checkBounds(c.Count * c.Count)
	if err := saveMathCount(db, c); err != nil {
		return nil, err
	}
	return c, nil
}
func Sqrt(db *mgo.Database) (*Counter, error) {
	c, err := LoadMathCount(db)
	if err != nil {
		return nil, err
	}
	c.Count = checkBounds(int(math.Sqrt(float64(c.Count))))
	if err := saveMathCount(db, c); err != nil {
		return nil, err
	}
	return c, nil
}

func checkBounds(count int) int {
	limit := 65536
	if count > limit {
		return limit
	} else if count < -limit {
		return -limit
	}
	return count
}

func LoadMathCount(db *mgo.Database) (*Counter, error) {
	mathdb := db.C("math")
	counter := &Counter{}
	query := bson.M{}
	if err := mathdb.Find(query).One(counter); err != nil {
		return counter, err
	}
	return counter, nil
}

func saveMathCount(db *mgo.Database, c *Counter) error {
	if err := db.C("math").Update(bson.M{}, c); err != nil {
		return err
	}
	return nil
}
