package models

import (
	"encoding/json"
	"fmt"

	"github.com/Sirupsen/logrus"
	"golang.org/x/net/websocket"
	"gopkg.in/mgo.v2"
)

type Client struct {
	WebSocket   *websocket.Conn
	WebSocketID string
	AccountID   string
	Username    string
}

var (
	WSConnect    = "WSCONNECT"
	WSDisconnect = "WSDISCONNECT"
)

var CommsActions = []string{
	WSConnect,
	WSDisconnect,
}

type CommsMessage struct {
	Message
}

type CommsQueue struct {
	DB      *mgo.Database
	Q       chan CommsMessage
	Clients map[string]Client //websocket id = Client
}

func NewCommsQueue(db *mgo.Database) (CommsQueue, error) {
	cq := CommsQueue{
		DB:      db,
		Clients: map[string]Client{},
	}

	cq.Q = make(chan CommsMessage)
	go cq.ReadMessages()
	return cq, nil
}

func (cq CommsQueue) ReadMessages() {
	for {
		commsMessage := <-cq.Q
		switch commsMessage.Type {
		case WSConnect:
			cq.HandleConnect(commsMessage)
		case WSDisconnect:
			cq.HandleDisconnect(cq.Clients[commsMessage.WebSocketID])
		default:
			continue
		}
	}
}

func (cq CommsQueue) HandleConnect(cm CommsMessage) {
	logrus.Infof("connecting %+v", cm.WebSocketID)
	if (cq.Clients[cm.WebSocketID] != Client{}) {
		return
	}
	cq.Clients[cm.WebSocketID] = Client{
		WebSocket:   cm.WebSocket,
		WebSocketID: cm.WebSocketID,
	}
}
func (cq CommsQueue) HandleDisconnect(client Client) {
	logrus.Infof("disconnecting %+v", client.WebSocketID)
	delete(cq.Clients, client.WebSocketID)
	client.WebSocket.Close()
}

func (cq CommsQueue) SendAll(msg interface{}) error {
	msgJSON, err := json.Marshal(msg)
	if err != nil {
		return fmt.Errorf("failed to sendall: %s", err.Error())
	}

	for _, client := range cq.Clients {
		if err := websocket.Message.Send(client.WebSocket, string(msgJSON)); err != nil {
			logrus.Warnf("discarding dead ws conn: %s", client.WebSocketID)
			cq.HandleDisconnect(client)
			continue
		}
	}
	return nil
}
