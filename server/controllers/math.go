package controllers

import (
	"fmt"
	"github.com/Sirupsen/logrus"
	"github.com/labstack/echo"
	"github.com/mikerjacobi/poker/server/models"
	"gopkg.in/mgo.v2"
)

type MathMessage struct {
	models.Message
	Count int `json:"count"`
}

func newMathMessage(action string, count int) MathMessage {
	return MathMessage{
		Message: models.Message{Type: action},
		Count:   count,
	}
}

func HandleIncrement(msg models.Message) error {
	db := msg.Context.Get("db").(*mgo.Database)

	count, err := models.IncrementCounter(db)
	if err != nil {
		models.SendError(msg.Sender.AccountID, "increment error")
		return fmt.Errorf("failed to increment: %+v", err)
	}

	if err := models.SendAll(newMathMessage(msg.Type, count)); err != nil {
		return fmt.Errorf("sendall error: %+v", err)
	}
	return nil
}

func HandleDecrement(msg models.Message) error {
	db := msg.Context.Get("db").(*mgo.Database)
	count, err := models.DecrementCounter(db)
	if err != nil {
		models.SendError(msg.Sender.AccountID, "decrement error")
		return fmt.Errorf("failed to decrement: %+v", err)
	}
	if err := models.SendAll(newMathMessage(msg.Type, count)); err != nil {
		return fmt.Errorf("sendall error: %+v", err)
	}
	return nil
}

func HandleSquare(msg models.Message) error {
	db := msg.Context.Get("db").(*mgo.Database)
	count, err := models.SquareCounter(db)
	if err != nil {
		models.SendError(msg.Sender.AccountID, "square error")
		return fmt.Errorf("failed to square: %+v", err)
	}
	if err := models.SendAll(newMathMessage(msg.Type, count)); err != nil {
		return fmt.Errorf("sendall error: %+v", err)
	}
	return nil
}

func HandleSqrt(msg models.Message) error {
	db := msg.Context.Get("db").(*mgo.Database)
	count, err := models.SqrtCounter(db)
	if err != nil {
		models.SendError(msg.Sender.AccountID, "sqrt error")
		return fmt.Errorf("failed to sqrt: %+v", err)
	}
	if err := models.SendAll(newMathMessage(msg.Type, count)); err != nil {
		return fmt.Errorf("sendall error: %+v", err)
	}
	return nil
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
