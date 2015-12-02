package models

import (
	"encoding/json"
	"fmt"

	"golang.org/x/net/websocket"
	"gopkg.in/mgo.v2"
)

type Client struct {
	WebSocket   *websocket.Conn
	WebSocketID string
	AccountID   string
	Username    string
}

type Message struct {
	Type string `json:"type"`
}

type MessageHandler struct {
	MathQueue
	Clients map[string]Client //websocket id = Client
}

func NewMessageHandler(db *mgo.Database) (MessageHandler, error) {
	mh := MessageHandler{}
	mq, err := NewMathQueue(db, &mh)
	if err != nil {
		return mh, fmt.Errorf("failed to init math queue: %s", err.Error())
	}

	mh.MathQueue = mq
	mh.Clients = map[string]Client{}
	return mh, nil
}

func (mh MessageHandler) Connect(wsID string, ws *websocket.Conn) {
	if (mh.Clients[wsID] != Client{}) {
		return
	}
	mh.Clients[wsID] = Client{
		WebSocket:   ws,
		WebSocketID: wsID,
	}
}

func (mh MessageHandler) Disconnect(wsID string) {
	delete(mh.Clients, wsID)
}

func (mh MessageHandler) Push(msg []byte, wsID string) error {
	//#TODO add wsID to mh.Clients

	m := Message{}
	if err := json.Unmarshal(msg, &m); err != nil {
		return err
	}

	if StringInSlice(m.Type, MathActions) {
		mathMessage := MathMessage{}
		if err := json.Unmarshal(msg, &mathMessage); err != nil {
			return err
		}
		mh.MathQueue.Q <- mathMessage
	} else {
		return fmt.Errorf("%s is an invalid action", m.Type)
	}
	return nil
}

func (mh MessageHandler) SendAll(msg interface{}) error {
	msgJSON, err := json.Marshal(msg)
	if err != nil {
		return fmt.Errorf("failed to sendall: %s", err.Error())
	}

	for _, client := range mh.Clients {
		if err := websocket.Message.Send(client.WebSocket, string(msgJSON)); err != nil {
			return err
		}
	}
	return nil
}
