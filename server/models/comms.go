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
	Account
}

type Comms struct {
	DB      *mgo.Database
	Clients map[string]*Client //websocket id = Client
}

func (c Comms) SetClient(client *Client) {
	if c.Clients[client.WebSocketID] != nil {
		//already connected
		return
	}
	c.Clients[client.WebSocketID] = client
}

func (c Comms) DeleteClient(wsID string) {
	if c.Clients[wsID] == nil {
		//doesn't exist, nothing to delete
		return
	}

	//close websocket
	c.Clients[wsID].WebSocket.Close()
	delete(c.Clients, wsID)
}

func (c Comms) Send(wsID string, msg interface{}) error {
	msgJSON, err := json.Marshal(msg)
	if err != nil {
		return fmt.Errorf("failed to jsonmarshal Comms.Send: %s", err.Error())
	}

	client := c.Clients[wsID]
	if client == nil {
		return fmt.Errorf("unknown websocket id in Comms.Send: %s", wsID)
	}

	if err := websocket.Message.Send(client.WebSocket, string(msgJSON)); err != nil {
		logrus.Warnf("discarding dead ws conn from send: %s", client.WebSocketID)
		c.DeleteClient(wsID)
	}
	return nil
}

func (c Comms) SendAll(msg interface{}) error {
	msgJSON, err := json.Marshal(msg)
	if err != nil {
		return fmt.Errorf("failed to sendall: %s", err.Error())
	}

	for _, client := range c.Clients {
		if err := websocket.Message.Send(client.WebSocket, string(msgJSON)); err != nil {
			logrus.Warnf("discarding dead ws conn from sendall: %s", client.WebSocketID)
			c.DeleteClient(client.WebSocketID)
		}
	}
	return nil
}
