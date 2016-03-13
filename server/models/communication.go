package models

import (
	"encoding/json"
	"fmt"

	"gopkg.in/mgo.v2"

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
		Type:  "/error",
		Error: errMsg,
	}

	if err := Send(accountID, msg); err != nil {
		logrus.Errorf("failed to SendError to %s: %+v", accountID, err)
	}
}

func SendAll(msg interface{}) error {
	for accountID := range connectionManager.Connections {
		if err := Send(accountID, msg); err != nil {
			logrus.Errorf("failed to SendAll to %s: %+v", accountID, err)
		}
	}
	return nil
}

func SendGame(db *mgo.Database, gameID string, msg interface{}) error {
	game, err := LoadGame(db, gameID, "")
	if err != nil {
		return fmt.Errorf("failed loadgame in SendGame: %+v", err)
	}

	for _, player := range game.Players {
		if err := Send(player.AccountID, msg); err != nil {
			logrus.Errorf("failed to SendGame to %s: %+v", player.AccountID, err)
		}
	}
	return nil
}

func SendGameError(db *mgo.Database, gameID string, errMsg string) {
	game, err := LoadGame(db, gameID, "")
	if err != nil {
		logrus.Errorf("failed to loadgame in sendGameError: %+v", err)
	}

	for _, player := range game.Players {
		SendError(player.AccountID, errMsg)
	}
}
