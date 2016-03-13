package models

import (
	"golang.org/x/net/websocket"
	"gopkg.in/mgo.v2"
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

func Connect(connection *Connection) {
	connectionManager.Connections[connection.AccountID] = connection
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
