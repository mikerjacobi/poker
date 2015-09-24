package controllers

import (
	"github.com/Sirupsen/logrus"
	"github.com/labstack/echo"
	"github.com/mikerjacobi/poker/models"
	"gopkg.in/mgo.v2"
)

func GetGame(c *echo.Context) error {
	logrus.Infof("in get game")
	db := c.Get("db").(*mgo.Database)
	gameID := c.Param("game_id")
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
	logrus.Infof("in get open games")
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

func validateCreateGame(gameName map[string][]string) error {
	return nil
}

func CreateGame(c *echo.Context) error {
	logrus.Infof("in create game")
	c.Request().ParseForm()
	err := validateCreateGame(c.Request().Form)
	if err != nil {
		logrus.Errorf("failed create account input validation %s", err.Error())
		c.JSON(400, Response{})
		return nil
	}

	gameName := c.Request().Form["gamename"][0]
	db := c.Get("db").(*mgo.Database)
	g, err := models.LoadGame(db, "", gameName)
	if err == nil {
		logrus.Errorf("game name taken: %s", g.Name)
		c.JSON(409, Response{})
		return nil
	} else if err != mgo.ErrNotFound && err != nil {
		logrus.Errorf("db error in create game: %s", err.Error())
		c.JSON(500, Response{})
		return nil
	}

	game, err := models.CreateGame(db, gameName)
	if err != nil {
		logrus.Errorf("failed to create game")
		c.JSON(500, Response{})
		return nil
	}
	c.JSON(200, Response{true, game})
	return nil
}

func JoinGame(c *echo.Context) error {
	gameID := c.Param("game_id")
	a, ok := c.Get("user").(models.Account)
	if !ok {
		logrus.Errorf("failed to get user in create game")
		c.JSON(500, Response{})
		return nil
	}

	logrus.Infof("game: %s, acctid: %s", gameID, a.AccountID)
	db := c.Get("db").(*mgo.Database)
	err := models.JoinGame(db, gameID, a.AccountID)
	if err == models.PlayerAlreadyJoined {
		logrus.Errorf("player already joined in join game")
		c.JSON(409, Response{})
		return nil
	} else if err != nil {
		logrus.Errorf("failed to join game")
		c.JSON(500, Response{})
		return nil
	}
	c.JSON(200, Response{Success: true})
	return nil
}
