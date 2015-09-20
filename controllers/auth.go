package controllers

import (
	"errors"

	"github.com/Sirupsen/logrus"
	"github.com/labstack/echo"
	"github.com/mikerjacobi/echomongo/models"
	"golang.org/x/crypto/bcrypt"
	"gopkg.in/mgo.v2"
)

func validateLogin(li map[string][]string) (models.Account, error) {
	a := models.Account{
		Username: li["username"][0],
		Password: li["password"][0],
	}
	if a.Username == "" {
		return a, errors.New("username cannot be empty")
	}
	if a.Password == "" {
		return a, errors.New("password cannot be empty")
	}
	return a, nil
}

func Login(c *echo.Context) error {
	logrus.Infof("login")

	c.Request().ParseForm()
	li := c.Request().Form
	inputAccount, err := validateLogin(li)
	if err != nil {
		logrus.Errorf("failed login input validation: %s", err.Error())
		c.JSON(400, Response{})
		return nil
	}

	db := c.Get("db").(*mgo.Database)
	a, err := models.LoadAccount(db, inputAccount.Username)
	if err != nil {
		logrus.Errorf("failed to load account in login")
		c.JSON(500, Response{})
		return err
	}

	if err := bcrypt.CompareHashAndPassword([]byte(a.Hashword), []byte(inputAccount.Password)); err != nil {
		logrus.Errorf("failed to authenticate in login: %s", err.Error())
		c.JSON(401, Response{})
		return nil
	}

	sessionID, err := a.NewSession(db)
	if err != nil {
		logrus.Errorf("failed to create new session in login: %s", err.Error())
		c.JSON(500, Response{})
		return nil
	}

	resp := struct {
		SessionID string `json:"session_id"`
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
