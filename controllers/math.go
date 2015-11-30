package controllers

import (
	"github.com/Sirupsen/logrus"
	"github.com/labstack/echo"
	"github.com/mikerjacobi/poker/models"
	"gopkg.in/mgo.v2"
)

func GetMathCount(c *echo.Context) error {
	logrus.Infof("in get math count")
	db := c.Get("db").(*mgo.Database)
	counter, err := models.LoadMathCount(db)
	if err != nil {
		logrus.Errorf("failed to get math count: %s", err.Error())
		c.JSON(500, Response{})
		return nil
	}
	c.JSON(200, Response{true, counter})
	return nil
}
