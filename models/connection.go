package models

import "gopkg.in/mgo.v2"

var (
	WSConnect    = "WSCONNECT"
	WSDisconnect = "WSDISCONNECT"
	WSError      = "WSERROR"
)

var ConnectionActions = []string{
	WSConnect,
	WSDisconnect,
}

type ConnectionMessage struct {
	Message
}

type ConnectionQueue struct {
	DB *mgo.Database
	Q  chan ConnectionMessage
	*Comms
}

func NewConnectionQueue(db *mgo.Database, comms *Comms) (ConnectionQueue, error) {
	cq := ConnectionQueue{
		DB:    db,
		Comms: comms,
	}

	cq.Q = make(chan ConnectionMessage)
	go cq.ReadMessages()
	return cq, nil
}

func (cq ConnectionQueue) ReadMessages() {
	for {
		cm := <-cq.Q
		switch cm.Type {
		case WSConnect:
			cq.SetClient(cm.Message)
		case WSDisconnect:
			cq.DeleteClient(cm.WebSocketID)
		default:
			continue
		}
	}
}
