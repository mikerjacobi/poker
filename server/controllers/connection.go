package controllers

import (
	"github.com/mikerjacobi/poker/server/models"
	"gopkg.in/mgo.v2"
)

var (
	WSConnect    = "WSCONNECT"
	WSDisconnect = "WSDISCONNECT"
)
var ConnectionActions = []string{
	WSConnect,
	WSDisconnect,
}

type ConnectionController struct {
	DB    *mgo.Database
	Queue chan Message
	*models.Comms
}

func newConnectionController(db *mgo.Database, comms *models.Comms) (ConnectionController, error) {
	cc := ConnectionController{
		DB:    db,
		Comms: comms,
	}

	cc.Queue = make(chan Message)
	go cc.ReadMessages()
	return cc, nil
}

func (cc ConnectionController) ReadMessages() {
	for {
		cm := <-cc.Queue
		switch cm.Type {
		case WSConnect:
			cc.SetClient(cm.Client())
		case WSDisconnect:
			cc.DeleteClient(cm.WebSocketID)
		default:
			continue
		}
	}
}
