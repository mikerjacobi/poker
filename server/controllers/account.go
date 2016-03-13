package controllers

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"

	"github.com/Sirupsen/logrus"
	"github.com/labstack/echo"
	"github.com/mikerjacobi/poker/server/models"
	"gopkg.in/mgo.v2"
)

const (
	accountLoad = "/account/load"
)

type CreateAccountRequest struct {
	LoginRequest
	Repeat string `json:"repeat"`
}

type AccountMessage struct {
	Type    string         `json:"type"`
	Account models.Account `json:"account"`
}

func newAccountMessage(action string, account models.Account) AccountMessage {
	return AccountMessage{
		Type:    action,
		Account: account,
	}
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

func HandleLoadAccount(msg models.Message) error {
	db, dbok := msg.Context.Get("db").(*mgo.Database)
	a, aok := msg.Context.Get("user").(models.Account)
	if !dbok || !aok {
		return fmt.Errorf("failed to load account or db.  dbok: %+v, aok: %+v", dbok, aok)
	}
	account, err := models.LoadAccount(db, a.Username)
	if err != nil {
		return fmt.Errorf("failed to load account in handle load account: %+v", err)
	}

	if err := models.Send(account.AccountID, newAccountMessage(accountLoad, account)); err != nil {
		return fmt.Errorf("failed to send in GetAccount: %+v", err)
	}
	return nil
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

func HandleChipRequest(msg models.Message) error {
	db, dbok := msg.Context.Get("db").(*mgo.Database)
	a, aok := msg.Context.Get("user").(models.Account)
	if !dbok || !aok {
		return fmt.Errorf("failed to load account or db.  dbok: %+v, aok: %+v", dbok, aok)
	}
	account, err := models.LoadAccount(db, a.Username)
	if err != nil {
		return fmt.Errorf("failed to load account in handle chip request: %+v", err)
	}

	req := struct {
		Amount int `json:"amount"`
	}{}
	if err := json.Unmarshal(msg.Raw, &req); err != nil {
		return fmt.Errorf("failed to unmarshal chipreq: %+v", err)
	}

	if req.Amount <= 0 {
		return fmt.Errorf("chip request <= 0: %+v", req)
	}
	account.Balance += req.Amount

	if err := account.Update(db); err != nil {
		return fmt.Errorf("failed to update account balance: %+v", err)
	}

	if err := models.Send(account.AccountID, newAccountMessage(accountLoad, account)); err != nil {
		return fmt.Errorf("failed to send in GetAccount: %+v", err)
	}
	return nil
}
