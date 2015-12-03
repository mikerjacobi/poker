package controllers

import (
	"errors"

	"github.com/Sirupsen/logrus"
	"github.com/labstack/echo"
	"github.com/mikerjacobi/poker/models"
	"golang.org/x/net/websocket"
	"gopkg.in/mgo.v2"
)

var mh models.MessageHandler

func InitializeMessageHandler(db *mgo.Database) error {
	var err error
	mh, err = models.NewMessageHandler(db)
	if err != nil {
		logrus.Errorf("failed to init message handler: %s", err.Error())
		return err
	}
	return nil
}

func HandleWebSocket(c *echo.Context) error {
	ws := c.Socket()
	msg := ""
	var wsID string

	for {
		if len(c.Request().Header["Sec-Websocket-Key"]) == 1 {
			wsID = c.Request().Header["Sec-Websocket-Key"][0]
		} else {
			e := "failed to pull websocket key from header"
			logrus.Errorf(e)
			return errors.New(e)
		}

		if err := websocket.Message.Receive(ws, &msg); err != nil {
			logrus.Errorf("failed to recv ws: %s", err.Error())
			return err
		}

		if err := mh.HandleMessage([]byte(msg), wsID, ws); err != nil {
			logrus.Errorf("failed to push msg %s because %s", msg, err.Error())
			continue
		}
	}
	return nil
}
