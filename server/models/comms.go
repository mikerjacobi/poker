package models

import (
	"encoding/json"
	"fmt"

	"github.com/Sirupsen/logrus"
	"github.com/labstack/echo"
	"golang.org/x/net/websocket"
	"gopkg.in/mgo.v2"
)

var (
	ServerError = "SERVERERROR"
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
type Message struct {
	Type            string          `json:"type"`
	WebSocketID     string          `json:"-"`
	WebSocket       *websocket.Conn `json:"-"`
	SenderAccountID string          `json:"-"`
	Raw             []byte          `json:"-"`
	Context         *echo.Context   `json:"-"`
}

type ErrorMessage struct {
	Type  string      `json:"type"`
	Error interface{} `json:"error"`
}

func (c Comms) SetClient(client *Client) {
	c.Clients[client.Account.AccountID] = client
}

func (c Comms) DeleteClient(accountID string) {
	_, ok := c.Clients[accountID]
	if !ok {
		//doesn't exist, nothing to delete
		return
	}

	//close websocket
	c.Clients[accountID].WebSocket.Close()
	delete(c.Clients, accountID)
}

func (c Comms) Send(accountID string, msg interface{}) error {
	msgJSON, err := json.Marshal(msg)
	if err != nil {
		return fmt.Errorf("failed to jsonmarshal Comms.Send: %s", err.Error())
	}

	client := c.Clients[accountID]
	if client == nil {
		return fmt.Errorf("unknown account id in Comms.Send: %s", accountID)
	}

	if err := websocket.Message.Send(client.WebSocket, string(msgJSON)); err != nil {
		logrus.Warnf("discarding dead ws conn from send: %s", client.WebSocketID)
		c.DeleteClient(accountID)
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

func (c Comms) SendGroup(msg interface{}, accountIDs []string) error {
	msgJSON, err := json.Marshal(msg)
	if err != nil {
		return fmt.Errorf("failed to sendgroupd: %s", err.Error())
	}

	for _, aid := range accountIDs {
		client, ok := c.Clients[aid]
		if !ok {
			return fmt.Errorf("accountID: %s is not currently connected", aid)
		}
		if err := websocket.Message.Send(client.WebSocket, string(msgJSON)); err != nil {
			logrus.Warnf("discarding dead ws conn from sendgroup: %s", client.WebSocketID)
			c.DeleteClient(client.WebSocketID)
		}
	}
	return nil
}
