package controllers

import (
	"github.com/Sirupsen/logrus"
	"github.com/mikerjacobi/poker/server/models"
	"gopkg.in/mgo.v2"
)

type ConnectionController struct {
	DB    *mgo.Database
	Queue chan models.Message
	*models.Comms
}

func newConnectionController(db *mgo.Database, comms *models.Comms) (ConnectionController, error) {
	cc := ConnectionController{
		DB:    db,
		Comms: comms,
	}

	cc.Queue = make(chan models.Message)
	go cc.ReadMessages()
	return cc, nil
}

func (cc ConnectionController) ReadMessages() {
	for {
		cm := <-cc.Queue
		switch cm.Type {
		case models.WSConnect:
			account, ok := cm.Context.Get("user").(models.Account)
			if !ok {
				logrus.Warnf("failed to get user in Client()")
				account = models.Account{}
			}
			cc.SetClient(&models.Client{
				WebSocket:   cm.WebSocket,
				WebSocketID: cm.WebSocketID,
				Account:     account,
			})
		case models.WSDisconnect:
			cc.DeleteClient(cm.WebSocketID)
		default:
			continue
		}
	}
}
