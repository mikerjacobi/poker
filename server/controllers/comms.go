package controllers

import (
	"github.com/Sirupsen/logrus"
	"github.com/mikerjacobi/poker/server/models"
	"golang.org/x/net/websocket"
	"gopkg.in/mgo.v2"
)

var (
	ServerError = "SERVERERROR"
)

type Message struct {
	Type        string          `json:"type"`
	WebSocketID string          `json:"-"`
	WebSocket   *websocket.Conn `json:"-"`
	Sender      models.Account  `json:"-"`
	Raw         []byte          `json:"-"`
}

func (m Message) Client() *models.Client {
	return &models.Client{
		WebSocket:   m.WebSocket,
		WebSocketID: m.WebSocketID,
		Account:     m.Sender,
	}
}

func newComms(db *mgo.Database) *models.Comms {
	c := models.Comms{}
	c.DB = db
	c.Clients = make(map[string]*models.Client)
	return &c
}

type ErrorMessage struct {
	Type  string      `json:"type"`
	Error interface{} `json:"error"`
}

func sendError(c *models.Comms, wsID string, msg interface{}) {
	err := ErrorMessage{
		Type:  ServerError,
		Error: msg,
	}
	if sendErr := c.Send(wsID, err); sendErr != nil {
		logrus.Errorf("send error: %+v", sendErr)
	}
}
