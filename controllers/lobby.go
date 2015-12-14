package controllers

import (
	"encoding/json"
	"io"
	"io/ioutil"

	"github.com/Sirupsen/logrus"
	"github.com/labstack/echo"
	"github.com/mikerjacobi/poker/models"
	"gopkg.in/mgo.v2"
)

type CreateGameRequest struct {
	GameName string `json:"gameName"`
	GameType string `json:"gameType"`
}

func GetGame(c *echo.Context) error {
	db := c.Get("db").(*mgo.Database)
	gameID := c.Param("gameID")
	game, err := models.LoadGame(db, gameID, "")
	if err != nil {
		logrus.Errorf("failed to get game: %s", err.Error())
		c.JSON(500, Response{})
		return nil
	}
	c.JSON(200, Response{true, game})
	return nil
}

func GetOpenGames(c *echo.Context) error {
	db := c.Get("db").(*mgo.Database)
	games, err := models.LoadOpenGames(db)
	if err != nil {
		logrus.Errorf("failed to get open games")
		c.JSON(500, Response{})
		return nil
	}
	c.JSON(200, Response{true, games})
	return nil
}

func validateCreateGame(cgBody io.ReadCloser) (*CreateGameRequest, error) {
	cgBytes, err := ioutil.ReadAll(cgBody)
	if err != nil {
		return nil, err
	}

	cg := CreateGameRequest{}
	err = json.Unmarshal(cgBytes, &cg)
	if err != nil {
		return nil, err
	}
	return &cg, nil
}

func CreateGame(c *echo.Context) error {
	cg, err := validateCreateGame(c.Request().Body)
	if err != nil {
		logrus.Errorf("failed create game input validation %s", err.Error())
		c.JSON(400, Response{})
		return nil
	}

	db := c.Get("db").(*mgo.Database)
	game, err := models.CreateGame(db, cg.GameName, cg.GameType)
	if err != nil {
		logrus.Errorf("failed to create game: %s", err)
		c.JSON(500, Response{})
		return nil
	}
	c.JSON(200, Response{true, game})
	return nil
}

func JoinGame(c *echo.Context) error {
	gameID := c.Param("gameID")
	a, ok := c.Get("user").(models.Account)
	if !ok {
		logrus.Errorf("failed to get user in create game")
		c.JSON(500, Response{})
		return nil
	}

	logrus.Infof("game: %s, acctid: %s", gameID, a.AccountID)
	db := c.Get("db").(*mgo.Database)
	game, err := models.JoinGame(db, gameID, a)
	if err == models.PlayerAlreadyJoined {
		logrus.Errorf("player already joined in join game")
		c.JSON(409, Response{})
		return nil
	} else if err != nil {
		logrus.Errorf("failed to join game")
		c.JSON(500, Response{})
		return nil
	}
	c.JSON(200, Response{Success: true, Payload: game})
	return nil
}
