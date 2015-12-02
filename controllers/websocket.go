package controllers

import (
	"fmt"

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
			logrus.Errorf("failed to pull websocket key from header")
			continue
		}

		if err := websocket.Message.Receive(ws, &msg); err != nil {
			if err.Error() == "EOF" {
				mh.Disconnect(wsID)
				continue
			} else {
				return fmt.Errorf("failed to recv ws: %s", err.Error())
			}
		}

		mh.Connect(wsID, ws)

		if err := mh.Push([]byte(msg), wsID); err != nil {
			logrus.Errorf("failed to push msg: %s", err.Error())
			continue
		}
		logrus.Infof("%+v: %s", wsID, msg)
	}
	return nil
}
