package controllers

import (
	"github.com/Sirupsen/logrus"
	"github.com/mikerjacobi/poker/server/models"
	"gopkg.in/mgo.v2"
)

func newComms(db *mgo.Database) *models.Comms {
	c := models.Comms{}
	c.DB = db
	c.Clients = make(map[string]*models.Client)
	return &c
}

func sendError(c *models.Comms, accountID string, msg interface{}) {
	err := models.ErrorMessage{
		Type:  models.ServerError,
		Error: msg,
	}
	if sendErr := c.Send(accountID, err); sendErr != nil {
		logrus.Errorf("send error: %+v", sendErr)
	}
}
