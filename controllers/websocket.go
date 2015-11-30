package controllers

import (
	"github.com/Sirupsen/logrus"
	"github.com/labstack/echo"
	"golang.org/x/net/websocket"
)

func HandleWebSocket(c *echo.Context) error {
	logrus.Infof("in handle websocket")
	ws := c.Socket()
	msg := ""

	for {
		if err := websocket.Message.Receive(ws, &msg); err != nil {
			return err
		}
		if err := websocket.Message.Send(ws, msg); err != nil {
			return err
		}
		logrus.Infof("%+v: %s", c.Request().Header["Sec-Websocket-Key"], msg)
	}
	return nil
}
