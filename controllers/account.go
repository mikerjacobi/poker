package controllers

import (
	"encoding/json"
	"errors"
	"io"
	"io/ioutil"

	"github.com/Sirupsen/logrus"
	"github.com/labstack/echo"
	"github.com/mikerjacobi/poker/models"
	"gopkg.in/mgo.v2"
)

type CreateAccountRequest struct {
	LoginRequest
	Repeat string `json:"repeat"`
}

func validateCreateAccount(carBody io.ReadCloser) (*CreateAccountRequest, error) {
	carBytes, err := ioutil.ReadAll(carBody)
	if err != nil {
		return nil, err
	}

	car := CreateAccountRequest{}
	err = json.Unmarshal(carBytes, &car)
	if err != nil {
		return nil, err
	}
	if car.Username == "" {
		return nil, errors.New("username cannot be empty")
	}
	if car.Password != car.Repeat {
		return nil, errors.New("passwords must match")
	}
	if car.Password == "" {
		return nil, errors.New("password cannot be empty")
	}
	return &car, nil
}

func CreateAccount(c *echo.Context) error {
	logrus.Infof("create account")
	caRequest, err := validateCreateAccount(c.Request().Body)
	if err != nil {
		logrus.Errorf("failed create account input validation %s", err.Error())
		c.JSON(400, Response{})
		return nil
	}
	db := c.Get("db").(*mgo.Database)
	_, err = models.LoadAccount(db, caRequest.Username)
	if err == nil {
		logrus.Errorf("account taken: %s", caRequest.Username)
		c.JSON(409, Response{})
		return nil
	} else if err != models.AccountNotFound && err != nil {
		logrus.Errorf("db error in create account: %s", err.Error())
		c.JSON(500, Response{})
		return nil
	}

	a := models.Account{
		Username: caRequest.Username,
		Password: caRequest.Password,
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
