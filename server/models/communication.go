package models

import (
	"encoding/json"
	"fmt"

	"github.com/Sirupsen/logrus"
	"github.com/labstack/echo"
	"golang.org/x/net/websocket"
)

type Message struct {
	Type       string        `json:"type"`
	Raw        []byte        `json:"-"`
	Context    *echo.Context `json:"-"`
	Sender     Account       `json:"-"`
	Connection *Connection   `json:"-"`
}

func Send(accountID string, msg interface{}) error {
	msgJSON, err := json.Marshal(msg)
	if err != nil {
		return fmt.Errorf("failed to jsonmarshal send: %s", err.Error())
	}

	conn, ok := connectionManager.Connections[accountID]
	if !ok {
		return fmt.Errorf("unknown accountID in send: %s", accountID)
	}

	if err := websocket.Message.Send(conn.Socket, string(msgJSON)); err != nil {
		logrus.Warnf("discarding dead ws conn from send for account: %s", conn.AccountID)
		Disconnect(accountID)
	}
	return nil
}

func SendError(accountID string, errMsg string) {
	msg := struct {
		Type  string `json:"type"`
		Error string `json:"error"`
	}{
		Type:  "SERVERERROR",
		Error: errMsg,
	}

	if err := Send(accountID, msg); err != nil {
		logrus.Errorf("failed to send error to %s: %+v", accountID, err)
	}
}

func SendAll(msg interface{}) error {
	msgJSON, err := json.Marshal(msg)
	if err != nil {
		return fmt.Errorf("failed to sendall: %s", err.Error())
	}

	for _, conn := range connectionManager.Connections {
		if err := websocket.Message.Send(conn.Socket, string(msgJSON)); err != nil {
			logrus.Warnf("discarding dead ws conn from sendall for account: %s", conn.AccountID)
			Disconnect(conn.AccountID)
		}
	}
	return nil
}

func SendGroup(accountIDs []string, msg interface{}) error {
	msgJSON, err := json.Marshal(msg)
	if err != nil {
		return fmt.Errorf("failed to sendgroupd: %s", err.Error())
	}

	for _, accountID := range accountIDs {
		conn, ok := connectionManager.Connections[accountID]
		if !ok {
			return fmt.Errorf("accountID: %s is not currently connected", accountID)
		}

		if err := websocket.Message.Send(conn.Socket, string(msgJSON)); err != nil {
			logrus.Warnf("discarding dead ws conn from sendgroup for account: %s", conn.AccountID)
			Disconnect(accountID)
		}
	}
	return nil
}
