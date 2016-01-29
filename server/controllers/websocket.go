package controllers

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/Sirupsen/logrus"
	"github.com/labstack/echo"
	"github.com/mikerjacobi/poker/server/models"
	"golang.org/x/net/websocket"
	"gopkg.in/mgo.v2"
)

var mh MessageHandler

type MessageHandler struct {
	*models.Comms
	MathController
	ConnectionController
	LobbyController
	HoldemController
	HighCardController
}

func InitializeMessageHandler(db *mgo.Database) error {

	comms := newComms(db)

	cc, err := newConnectionController(db, comms)
	if err != nil {
		return fmt.Errorf("failed to init connection controller: %s", err.Error())
	}

	mc, err := newMathController(db, comms)
	if err != nil {
		return fmt.Errorf("failed to init math controller: %s", err.Error())
	}

	hc, err := newHoldemController(db, comms)
	if err != nil {
		return fmt.Errorf("failed to init holdem controller: %s", err.Error())
	}

	hcc, err := newHighCardController(db, comms)
	if err != nil {
		return fmt.Errorf("failed to init highcard controller: %s", err.Error())
	}

	lc, err := newLobbyController(db, comms, hc, hcc)
	if err != nil {
		return fmt.Errorf("failed to init lobby controller: %s", err.Error())
	}

	mh = MessageHandler{comms, mc, cc, lc, hc, hcc}
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

		account, ok := c.Get("user").(models.Account)
		if !ok {
			logrus.Errorf("failed to get user in handle websocket")
			continue
		}

		if err := websocket.Message.Receive(ws, &msg); err != nil {
			//close connection gracefully
			return c.JSON(200, Response{true, nil})
		}

		if err := mh.HandleMessage(c, []byte(msg), wsID, ws, account); err != nil {
			logrus.Errorf("failed to push msg %s: %s", msg, err.Error())
			continue
		}
	}
	return nil
}

func (mh MessageHandler) HandleMessage(c *echo.Context, msg []byte, wsID string, ws *websocket.Conn, a models.Account) error {
	m := models.Message{
		WebSocketID:     wsID,
		WebSocket:       ws,
		SenderAccountID: a.AccountID,
		Raw:             msg,
		Context:         c,
	}
	if err := json.Unmarshal(msg, &m); err != nil {
		return err
	}

	logrus.Infof("%s: %+s", a.Username, m.Type)
	//logrus.Infof("%s: %+s %+v", m.Sender.Username, m.Type, string(msg))
	if models.StringInSlice(m.Type, MathActions) {
		mh.MathController.Queue <- m
	} else if models.StringInSlice(m.Type, models.ConnectionActions) {
		mh.ConnectionController.Queue <- m
	} else if models.StringInSlice(m.Type, models.LobbyActions) {
		mh.LobbyController.Queue <- m
	} else if models.StringInSlice(m.Type, models.HoldemActions) {
		mh.HoldemController.Queue <- m
	} else if models.StringInSlice(m.Type, models.HighCardActions) {
		mh.HighCardController.Queue <- m
	} else {
		err := fmt.Sprintf("invalid action: %s", m.Type)
		logrus.Errorf("%s: %s", a.AccountID, err)
		sendError(mh.ConnectionController.Comms, a.AccountID, err)
	}
	return nil
}
