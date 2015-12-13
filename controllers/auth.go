package controllers

import (
	"encoding/json"
	"errors"
	"io"
	"io/ioutil"

	"github.com/Sirupsen/logrus"
	"github.com/labstack/echo"
	"github.com/mikerjacobi/poker/models"
	"golang.org/x/crypto/bcrypt"
	"gopkg.in/mgo.v2"
)

type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func validateLogin(lrBody io.ReadCloser) (*LoginRequest, error) {
	lrBytes, err := ioutil.ReadAll(lrBody)
	if err != nil {
		return nil, err
	}

	lr := LoginRequest{}
	err = json.Unmarshal(lrBytes, &lr)
	if err != nil {
		return nil, err
	}

	if lr.Username == "" {
		return nil, errors.New("username cannot be empty")
	}

	if lr.Password == "" {
		return nil, errors.New("password cannot be empty")
	}
	return &lr, nil
}

func Login(c *echo.Context) error {
	logrus.Infof("login")

	loginRequest, err := validateLogin(c.Request().Body)
	if err != nil {
		logrus.Errorf("failed login validation: %s", err.Error())
		c.JSON(400, Response{})
		return nil
	}

	db := c.Get("db").(*mgo.Database)
	account, err := models.LoadAccount(db, loginRequest.Username)
	if err != nil {
		logrus.Errorf("failed to load account in login: %s", err)
		c.JSON(500, Response{})
		return nil
	}

	if err := bcrypt.CompareHashAndPassword([]byte(account.Hashword), []byte(loginRequest.Password)); err != nil {
		logrus.Errorf("failed to authenticate in login: %s", err.Error())
		c.JSON(401, Response{})
		return nil
	}

	sessionID, err := account.NewSession(db)
	if err != nil {
		logrus.Errorf("failed to create new session in login: %s", err.Error())
		c.JSON(500, Response{})
		return nil
	}

	resp := struct {
		SessionID string `json:"sessionID"`
	}{sessionID}

	c.JSON(200, Response{
		Success: true,
		Payload: resp,
	})
	return nil
}

func Logout(c *echo.Context) error {
	logrus.Infof("logout")

	a, ok := c.Get("user").(models.Account)
	if !ok {
		logrus.Errorf("failed to get user in logout")
		c.JSON(500, Response{})
		return nil
	}

	db := c.Get("db").(*mgo.Database)
	err := a.ClearSession(db)
	if err != nil {
		logrus.Errorf("failed to clear session in logout: %s", err.Error())
		c.JSON(500, Response{})
		return nil
	}

	c.JSON(200, Response{Success: true})
	return nil
}
