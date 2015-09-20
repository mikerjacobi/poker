package controllers

import (
	"errors"

	"github.com/Sirupsen/logrus"
	"github.com/labstack/echo"
	"github.com/mikerjacobi/echomongo/models"
	"gopkg.in/mgo.v2"
)

func validateCreateAccount(cai map[string][]string) (models.Account, error) {
	a := models.Account{
		Username: cai["username"][0],
		Password: cai["pw1"][0],
	}
	if a.Username == "" {
		return a, errors.New("username cannot be empty")
	}
	if a.Password != cai["pw2"][0] {
		return a, errors.New("passwords must match")
	}
	if a.Password == "" {
		return a, errors.New("password cannot be empty")
	}
	return a, nil
}

func CreateAccount(c *echo.Context) error {
	logrus.Infof("create account")
	c.Request().ParseForm()
	cai := c.Request().Form
	a, err := validateCreateAccount(cai)
	if err != nil {
		logrus.Errorf("failed create account input validation %s", err.Error())
		c.JSON(400, Response{})
		return nil
	}
	db := c.Get("db").(*mgo.Database)
	_, err = models.LoadAccount(db, a.Username)
	if err == nil {
		logrus.Errorf("account taken: %s", a.Username)
		c.JSON(409, Response{})
		return nil
	} else if err != models.AccountNotFound && err != nil {
		logrus.Errorf("db error in create account: %s", err.Error())
		c.JSON(500, Response{})
		return nil
	}

	err = models.CreateAccount(db, a)
	if err != nil {
		logrus.Errorf("failed to create account: %s", err.Error())
		c.JSON(500, Response{})
		return nil
	}

	c.JSON(200, Response{true, a})
	return nil
}
