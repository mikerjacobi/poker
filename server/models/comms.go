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
	//check to see if this wsID/accountid is connected in another session
	_, acctOK := c.Clients[client.Account.AccountID]
	_, wsOK := c.Clients[client.WebSocketID]
	if acctOK && wsOK {
		//already connected correctly
		return
	} else if acctOK || wsOK {
		//this means account and websocket are out of sync
		//delete the bad connection and continue
		c.DeleteClient(client.WebSocketID)
	}

	//map both websocketID and accountID to the same client; allows for easier access later
	c.Clients[client.WebSocketID] = client
	c.Clients[client.Account.AccountID] = client
}

func (c Comms) DeleteClient(wsID string) {
	_, wsOK := c.Clients[wsID]
	if !wsOK {
		//doesn't exist, nothing to delete
		return
	}

	//close websocket
	c.Clients[wsID].WebSocket.Close()

	//remove both references to this connection
	delete(c.Clients, c.Clients[wsID].Account.AccountID)
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
