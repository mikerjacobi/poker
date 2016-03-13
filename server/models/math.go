package models

import (
	"fmt"
	"math"

	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

func IncrementCounter(db *mgo.Database) (int, error) {
	c, err := LoadMathCount(db)
	if err != nil {
		return 0, err
	}
	c.Count = checkBounds(c.Count + 1)
	if err := saveMathCount(db, c); err != nil {
		return 0, err
	}
	return c.Count, nil
}
func DecrementCounter(db *mgo.Database) (int, error) {
	c, err := LoadMathCount(db)
	if err != nil {
		return 0, err
	}
	c.Count = checkBounds(c.Count - 1)
	if err := saveMathCount(db, c); err != nil {
		return 0, err
	}
	return c.Count, nil
}
func SquareCounter(db *mgo.Database) (int, error) {
	c, err := LoadMathCount(db)
	if err != nil {
		return 0, err
	}
	c.Count = checkBounds(c.Count * c.Count)
	if err := saveMathCount(db, c); err != nil {
		return 0, err
	}
	return c.Count, nil
}
func SqrtCounter(db *mgo.Database) (int, error) {
	c, err := LoadMathCount(db)
	if err != nil {
		return 0, err
	}
	c.Count = checkBounds(int(math.Sqrt(float64(c.Count))))
	if err := saveMathCount(db, c); err != nil {
		return 0, err
	}
	return c.Count, nil
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

type Counter struct {
	Count int `json:"count" bson:"count"`
}

func LoadMathCount(db *mgo.Database) (*Counter, error) {
	mathdb := db.C("math")
	counter := Counter{}
	if err := mathdb.Find(bson.M{}).One(&counter); err != nil {
		if err.Error() == "not found" {
			//look for uninitialized math counter...
			if err := initializeMathCounter(db); err != nil {
				return nil, err
			}
			counter.Count = 0
		} else {
			return nil, err
		}
	}
	return &counter, nil
}

func initializeMathCounter(db *mgo.Database) error {
	mathDB := db.C("math")
	//we have an unintialized math counter
	if err := mathDB.Insert(Counter{0}); err != nil {
		return fmt.Errorf("failed to init math counter: %s", err)
	}
	return nil
}

func saveMathCount(db *mgo.Database, c *Counter) error {
	if err := db.C("math").Update(bson.M{}, c); err != nil {
		return err
	}
	return nil
}
