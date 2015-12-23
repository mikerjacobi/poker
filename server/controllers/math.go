package controllers

import (
	"github.com/Sirupsen/logrus"
	"github.com/labstack/echo"
	"github.com/mikerjacobi/poker/server/models"
	"gopkg.in/mgo.v2"
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

type MathMessage struct {
	Message
	*models.Counter
}

type MathController struct {
	DB    *mgo.Database
	Queue chan Message
	*models.Comms
}

func newMathController(db *mgo.Database, c *models.Comms) (MathController, error) {
	mc := MathController{
		DB:    db,
		Comms: c,
	}

	mc.Queue = make(chan Message)
	go mc.ReadMessages()
	return mc, nil
}

func (mc MathController) ReadMessages() {
	for {
		m := <-mc.Queue
		switch m.Type {
		case Increment:
			mc.HandleIncrement(m)
		case Decrement:
			mc.HandleDecrement(m)
		case Square:
			mc.HandleSquare(m)
		case Sqrt:
			mc.HandleSqrt(m)
		default:
			continue
		}
	}
}

func (mc MathController) HandleIncrement(msg Message) {
	log := logrus.WithFields(logrus.Fields{"func": "HandleIncrement"})
	c, err := models.Increment(mc.DB)
	if err != nil {
		e := "increment error "
		sendError(mc.Comms, msg.WebSocketID, e)
		logrus.Errorf("%s: %s", msg.Sender.AccountID, e+err.Error())
		return
	}
	incrementMsg := MathMessage{Message: msg, Counter: c}
	if err := mc.SendAll(incrementMsg); err != nil {
		log.Errorf("sendall error: %+v", err)
		return
	}
}

func (mc MathController) HandleDecrement(msg Message) {
	log := logrus.WithFields(logrus.Fields{"func": "HandleDecrement"})
	c, err := models.Decrement(mc.DB)
	if err != nil {
		e := "decrement error "
		sendError(mc.Comms, msg.WebSocketID, e)
		logrus.Errorf("%s: %s", msg.Sender.AccountID, e+err.Error())
		return
	}
	decrementMsg := MathMessage{Message: msg, Counter: c}
	if err := mc.SendAll(decrementMsg); err != nil {
		log.Errorf("sendall error: %+v", err)
		return
	}
}

func (mc MathController) HandleSquare(msg Message) {
	log := logrus.WithFields(logrus.Fields{"func": "HandleSquare"})
	c, err := models.Square(mc.DB)
	if err != nil {
		e := "square error "
		sendError(mc.Comms, msg.WebSocketID, e)
		logrus.Errorf("%s: %s", msg.Sender.AccountID, e+err.Error())
		return
	}
	squareMsg := MathMessage{Message: msg, Counter: c}
	if err := mc.SendAll(squareMsg); err != nil {
		log.Errorf("sendall error: %+v", err)
		return
	}
}

func (mc MathController) HandleSqrt(msg Message) {
	log := logrus.WithFields(logrus.Fields{"func": "HandleSquare"})
	c, err := models.Sqrt(mc.DB)
	if err != nil {
		e := "sqrt error "
		sendError(mc.Comms, msg.WebSocketID, e)
		logrus.Errorf("%s: %s", msg.Sender.AccountID, e+err.Error())
		return
	}
	squareMsg := MathMessage{Message: msg, Counter: c}
	if err := mc.SendAll(squareMsg); err != nil {
		log.Errorf("sendall error: %+v", err)
		return
	}
}

func GetMathCount(c *echo.Context) error {
	db := c.Get("db").(*mgo.Database)
	counter, err := models.LoadMathCount(db)
	if err != nil {
		logrus.Errorf("failed to get math count: %s", err.Error())
		c.JSON(500, Response{})
		return nil
	}
	c.JSON(200, Response{true, counter})
	return nil
}
