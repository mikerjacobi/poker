package controllers

import (
	"github.com/Sirupsen/logrus"
	"github.com/labstack/echo"
	"github.com/mikerjacobi/poker/server/models"
	"gopkg.in/mgo.v2"
)

var MathActions = []string{
	models.Increment,
	models.Decrement,
	models.Square,
	models.Sqrt,
}

type MathController struct {
	DB    *mgo.Database
	Queue chan models.Message
	*models.Comms
}

func newMathController(db *mgo.Database, c *models.Comms) (MathController, error) {
	mc := MathController{
		DB:    db,
		Comms: c,
	}

	mc.Queue = make(chan models.Message)
	go mc.ReadMessages()
	return mc, nil
}

func (mc MathController) ReadMessages() {
	for {
		m := <-mc.Queue
		switch m.Type {
		case models.Increment:
			mc.HandleIncrement(m)
		case models.Decrement:
			mc.HandleDecrement(m)
		case models.Square:
			mc.HandleSquare(m)
		case models.Sqrt:
			mc.HandleSqrt(m)
		default:
			continue
		}
	}
}

func (mc MathController) HandleIncrement(msg models.Message) {
	log := logrus.WithFields(logrus.Fields{"func": "HandleIncrement"})
	c, err := models.IncrementCounter(mc.DB)
	if err != nil {
		e := "increment error "
		sendError(mc.Comms, msg.SenderAccountID, e)
		logrus.Errorf("%s: %s", msg.SenderAccountID, e+err.Error())
		return
	}
	incrementMsg := models.MathMessage{Message: msg, Counter: c}
	if err := mc.SendAll(incrementMsg); err != nil {
		log.Errorf("sendall error: %+v", err)
		return
	}
}

func (mc MathController) HandleDecrement(msg models.Message) {
	log := logrus.WithFields(logrus.Fields{"func": "HandleDecrement"})
	c, err := models.DecrementCounter(mc.DB)
	if err != nil {
		e := "decrement error "
		sendError(mc.Comms, msg.SenderAccountID, e)
		logrus.Errorf("%s: %s", msg.SenderAccountID, e+err.Error())
		return
	}
	decrementMsg := models.MathMessage{Message: msg, Counter: c}
	if err := mc.SendAll(decrementMsg); err != nil {
		log.Errorf("sendall error: %+v", err)
		return
	}
}

func (mc MathController) HandleSquare(msg models.Message) {
	log := logrus.WithFields(logrus.Fields{"func": "HandleSquare"})
	c, err := models.SquareCounter(mc.DB)
	if err != nil {
		e := "square error "
		sendError(mc.Comms, msg.SenderAccountID, e)
		logrus.Errorf("%s: %s", msg.SenderAccountID, e+err.Error())
		return
	}
	squareMsg := models.MathMessage{Message: msg, Counter: c}
	if err := mc.SendAll(squareMsg); err != nil {
		log.Errorf("sendall error: %+v", err)
		return
	}
}

func (mc MathController) HandleSqrt(msg models.Message) {
	log := logrus.WithFields(logrus.Fields{"func": "HandleSquare"})
	c, err := models.SqrtCounter(mc.DB)
	if err != nil {
		e := "sqrt error "
		sendError(mc.Comms, msg.SenderAccountID, e)
		logrus.Errorf("%s: %s", msg.SenderAccountID, e+err.Error())
		return
	}
	squareMsg := models.MathMessage{Message: msg, Counter: c}
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
