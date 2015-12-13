package models

import (
	"encoding/json"
	"fmt"

	"github.com/Sirupsen/logrus"
	"golang.org/x/net/websocket"
	"gopkg.in/mgo.v2"
)

type Message struct {
	Type        string          `json:"type"`
	WebSocketID string          `json:"-"`
	WebSocket   *websocket.Conn `json:"-"`
	Sender      Account         `json:"-"`
}

type ErrorMessage struct {
	Type  string      `json:"type"`
	Error interface{} `json:"error"`
}

type MessageHandler struct {
	Comms
	MathQueue
	ConnectionQueue
	LobbyQueue
	HoldemQueue
}

func newErrorMessage(err interface{}) ErrorMessage {
	return ErrorMessage{
		Type:  WSError,
		Error: err,
	}
}

func NewMessageHandler(db *mgo.Database) (MessageHandler, error) {

	comms := newComms(db)

	cq, err := NewConnectionQueue(db, &comms)
	if err != nil {
		return MessageHandler{}, fmt.Errorf("failed to init connection queue: %s", err.Error())
	}

	mq, err := NewMathQueue(db, &comms)
	if err != nil {
		return MessageHandler{}, fmt.Errorf("failed to init math queue: %s", err.Error())
	}

	lq, err := NewLobbyQueue(db, &comms)
	if err != nil {
		return MessageHandler{}, fmt.Errorf("failed to init lobby queue: %s", err.Error())
	}

	hq, err := NewHoldemQueue(db, &comms)
	if err != nil {
		return MessageHandler{}, fmt.Errorf("failed to init holdem queue: %s", err.Error())
	}

	mh := MessageHandler{comms, mq, cq, lq, hq}
	return mh, nil
}

func (mh MessageHandler) HandleMessage(msg []byte, wsID string, ws *websocket.Conn, a Account) error {
	m := Message{
		WebSocketID: wsID,
		WebSocket:   ws,
		Sender:      a,
	}
	if err := json.Unmarshal(msg, &m); err != nil {
		return err
	}

	logrus.Infof("%s: %+s", m.Sender.Username, m.Type)
	//logrus.Infof("%s: %+s %+v", m.Sender.Username, m.Type, msg)
	if StringInSlice(m.Type, MathActions) {
		mathMessage := MathMessage{Message: m}
		if err := json.Unmarshal(msg, &mathMessage); err != nil {
			return err
		}
		mh.MathQueue.Q <- mathMessage
	} else if StringInSlice(m.Type, ConnectionActions) {
		connectionMessage := ConnectionMessage{Message: m}
		if err := json.Unmarshal(msg, &connectionMessage); err != nil {
			return err
		}
		mh.ConnectionQueue.Q <- connectionMessage
	} else if StringInSlice(m.Type, LobbyActions) {
		lobbyMessage := LobbyMessage{Message: m}
		if err := json.Unmarshal(msg, &lobbyMessage); err != nil {
			return err
		}
		mh.LobbyQueue.Q <- lobbyMessage
	} else if StringInSlice(m.Type, HoldemActions) {
		holdemMessage := HoldemMessage{Message: m}
		if err := json.Unmarshal(msg, &holdemMessage); err != nil {
			return err
		}
		mh.HoldemQueue.Q <- holdemMessage
	} else {
		err := struct {
			InvalidAction string `json:"invalid_action"`
		}{m.Type}
		wsError := newErrorMessage(err)
		mh.ConnectionQueue.Send(wsID, wsError)
	}
	return nil
}
