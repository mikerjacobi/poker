package models

import (
	"math"

	"github.com/Sirupsen/logrus"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

var (
	Increment = "INCREMENT"
	Decrement = "DECREMENT"
	Square    = "SQUARE"
	Sqrt      = "SQRT"
)

var MathActions = []string{
	Increment,
	Decrement,
	Square,
	Sqrt,
}

type Counter struct {
	Count int `json:"count" bson:"count"`
}

type MathMessage struct {
	Message
	Counter
}

type MathQueue struct {
	Count int
	DB    *mgo.Database
	Q     chan MathMessage
	//SendAll func(interface{}) error
	MH *MessageHandler
}

func NewMathQueue(db *mgo.Database, mh *MessageHandler) (MathQueue, error) {
	mq := MathQueue{
		DB: db,
		//SendAll: sendAll,
		MH: mh,
	}

	c, err := LoadMathCount(db)
	if err != nil {
		return mq, err
	}
	mq.Count = c.Count
	mq.Q = make(chan MathMessage)
	go mq.ReadMessages()
	return mq, nil
}

func (mq MathQueue) ReadMessages() {
	for {
		mathMessage := <-mq.Q
		switch mathMessage.Type {
		case Increment:
			if err := mq.HandleIncrement(); err != nil {
				logrus.Errorf("failed to increment")
				continue
			}
		case Decrement:
			if err := mq.HandleDecrement(); err != nil {
				logrus.Errorf("failed to decrement")
				continue
			}
		case Square:
			if err := mq.HandleSquare(); err != nil {
				logrus.Errorf("failed to square")
				continue
			}
		case Sqrt:
			if err := mq.HandleSqrt(); err != nil {
				logrus.Errorf("failed to sqrt")
				continue
			}
		default:
			continue
		}
	}
}

func (mq MathQueue) HandleIncrement() error {
	c, err := LoadMathCount(mq.DB)
	if err != nil {
		return err
	}
	c.Count = checkBounds(c.Count + 1)
	if err := saveMathCount(mq.DB, c); err != nil {
		return err
	}
	m := MathMessage{
		Message: Message{Type: Increment},
		Counter: c,
	}
	return mq.MH.SendAll(m)
}
func (mq MathQueue) HandleDecrement() error {
	c, err := LoadMathCount(mq.DB)
	if err != nil {
		return err
	}
	c.Count = checkBounds(c.Count - 1)
	if err := saveMathCount(mq.DB, c); err != nil {
		return err
	}
	m := MathMessage{
		Message: Message{Type: Increment},
		Counter: c,
	}
	return mq.MH.SendAll(m)
}
func (mq MathQueue) HandleSquare() error {
	c, err := LoadMathCount(mq.DB)
	if err != nil {
		return err
	}
	c.Count = checkBounds(c.Count * c.Count)
	if err := saveMathCount(mq.DB, c); err != nil {
		return err
	}
	m := MathMessage{
		Message: Message{Type: Increment},
		Counter: c,
	}
	return mq.MH.SendAll(m)
}
func (mq MathQueue) HandleSqrt() error {
	c, err := LoadMathCount(mq.DB)
	if err != nil {
		return err
	}
	c.Count = checkBounds(int(math.Sqrt(float64(c.Count))))
	if err := saveMathCount(mq.DB, c); err != nil {
		return err
	}
	m := MathMessage{
		Message: Message{Type: Increment},
		Counter: c,
	}
	return mq.MH.SendAll(m)
}

func checkBounds(count int) int {
	if count > 65536 {
		return 65536
	}
	if count < -65536 {
		return -65536
	}
	return count
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

func saveMathCount(db *mgo.Database, c Counter) error {
	mathdb := db.C("math")
	if err := mathdb.Update(bson.M{}, c); err != nil {
		return err
	}
	return nil
}
