package models

import (
	"golang.org/x/net/websocket"
	"gopkg.in/mgo.v2"
)

const (
	wsInfo = "/ws/info"
)

var connectionManager *ConnectionManager

type Connection struct {
	AccountID string
	SocketID  string
	Socket    *websocket.Conn
}

type ConnectionManager struct {
	DB          *mgo.Database
	Connections map[string]*Connection //account id = Connection
}

func InitializeConnectionManager(db *mgo.Database) {
	cm := ConnectionManager{
		DB:          db,
		Connections: make(map[string]*Connection),
	}
	//TODO once.Do this
	connectionManager = &cm
}

func Connect(connection *Connection, account Account) {
	connectionManager.Connections[connection.AccountID] = connection
	msg := struct {
		Type      string `json:"type"`
		AccountID string `json:"accountID"`
		Username  string `json:"username"`
	}{wsInfo, account.AccountID, account.Username}
	Send(account.AccountID, msg)
}

func Disconnect(accountID string) {
	_, ok := connectionManager.Connections[accountID]
	if !ok {
		//doesn't exist, nothing to delete
		return
	}

	//close websocket
	connectionManager.Connections[accountID].Socket.Close()
	delete(connectionManager.Connections, accountID)
}
