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
	Handlers map[string]func(models.Message) error
}

func InitializeMessageHandler(db *mgo.Database) (*MessageHandler, error) {

	models.InitializeConnectionManager(db)
	models.InitializeHighCardManager(db)
	mh := MessageHandler{Handlers: map[string]func(models.Message) error{}}
	return &mh, nil
}

func (mh MessageHandler) Handle(action string, actionHandler func(models.Message) error) {
	mh.Handlers[action] = actionHandler
}

func (mh *MessageHandler) HandleWebSocket(c *echo.Context) error {
	ws := c.Socket()
	rawMsg := ""
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

		if err := websocket.Message.Receive(ws, &rawMsg); err != nil {
			//close connection gracefully
			return c.JSON(200, Response{true, nil})
		}

		msg := models.Message{
			Connection: &models.Connection{account.AccountID, wsID, ws},
			Sender:     account,
			Raw:        []byte(rawMsg),
			Context:    c,
		}
		if err := json.Unmarshal(msg.Raw, &msg); err != nil {
			logrus.Errorf("failed to unmarshal msg %s: %s", rawMsg, err.Error())
			continue
		}

		handleAction := mh.Handlers[msg.Type]
		if handleAction == nil {
			msg.Type = "/default"
			handleAction = mh.Handlers["/default"]
		}

		logrus.Infof("%s: %+s", msg.Sender.Username, msg.Type)
		if err := handleAction(msg); err != nil {
			logrus.Errorf("failed to handle %s's action: %s.  %+v", msg.Sender.Username, msg.Type, err)
		}
	}
	return nil
}

func DefaultActionHandler(msg models.Message) error {
	return fmt.Errorf("invalid payload: %+v", string(msg.Raw))
}

func HandleWebSocketConnect(msg models.Message) error {
	account, ok := msg.Context.Get("user").(models.Account)
	if !ok {
		return fmt.Errorf("failed to pull user out of conext in handle ws connect")
	}
	if account.AccountID != msg.Connection.AccountID {
		return fmt.Errorf("accountid mismatch in handle ws connect")
	}
	models.Connect(msg.Connection, account)
	return nil
}

func HandleWebSocketDisconnect(msg models.Message) error {
	models.Disconnect(msg.Sender.AccountID)
	return nil
}
