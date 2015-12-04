package models

import (
	"encoding/json"
	"fmt"

	"github.com/Sirupsen/logrus"
	"golang.org/x/net/websocket"
	"gopkg.in/mgo.v2"
)

type Message struct {
	Type        string `json:"type"`
	WebSocketID string
	WebSocket   *websocket.Conn
}

type MessageHandler struct {
	MathQueue
	CommsQueue
}

func NewMessageHandler(db *mgo.Database) (MessageHandler, error) {

	cq, err := NewCommsQueue(db)
	if err != nil {
		return MessageHandler{}, fmt.Errorf("failed to init comms queue: %s", err.Error())
	}

	mq, err := NewMathQueue(db, &cq)
	if err != nil {
		return MessageHandler{}, fmt.Errorf("failed to init math queue: %s", err.Error())
	}
	mh := MessageHandler{mq, cq}

	return mh, nil
}

func (mh MessageHandler) HandleMessage(msg []byte, wsID string, ws *websocket.Conn) error {
	m := Message{
		WebSocketID: wsID,
		WebSocket:   ws,
	}
	if err := json.Unmarshal(msg, &m); err != nil {
		return err
	}

	if StringInSlice(m.Type, MathActions) {
		mathMessage := MathMessage{Message: m}
		if err := json.Unmarshal(msg, &mathMessage); err != nil {
			return err
		}
		logrus.Infof("%+v", mathMessage)
		mh.MathQueue.Q <- mathMessage
	} else if StringInSlice(m.Type, CommsActions) {
		commsMessage := CommsMessage{Message: m}
		if err := json.Unmarshal(msg, &commsMessage); err != nil {
			return err
		}
		logrus.Infof("%+v", commsMessage)
		mh.CommsQueue.Q <- commsMessage
	} else {
		return fmt.Errorf("%s is an invalid action", m.Type)
	}
	return nil
}
