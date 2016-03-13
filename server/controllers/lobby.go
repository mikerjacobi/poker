package controllers

import (
	"encoding/json"
	"fmt"

	"github.com/Sirupsen/logrus"
	"github.com/labstack/echo"
	"github.com/mikerjacobi/poker/server/models"
	"gopkg.in/mgo.v2"
)

type LobbyMessage struct {
	Type        string `json:"type"`
	models.Game `json:"game"`
}

func newLobbyMessage(action string, game models.Game) LobbyMessage {
	return LobbyMessage{
		Type: action,
		Game: game,
	}
}

var (
	GameTypes []string
)

type CreateGameRequest struct {
	Name string `json:"gameName"`
	Type string `json:"gameType"`
}
type JoinLeaveGameRequest struct {
	ID string `json:"gameID"`
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

func validateCreateGame(msg models.Message) (*CreateGameRequest, error) {
	cg := struct {
		Game CreateGameRequest `json:"game"`
	}{}
	err := json.Unmarshal(msg.Raw, &cg)
	if err != nil {
		return nil, err
	}

	if cg.Game.Name == "" {
		return nil, fmt.Errorf("gamename cannot be empty: %+v", string(msg.Raw))
	}

	if !models.StringInSlice(cg.Game.Type, []string{"highcard"}) {
		return nil, fmt.Errorf("invalid gametype: %s", cg.Game.Type)
	}
	return &cg.Game, nil
}

func HandleCreateGame(msg models.Message) error {
	userFailMsg := "failed to create game"
	cg, err := validateCreateGame(msg)
	if err != nil {
		models.SendError(msg.Sender.AccountID, userFailMsg)
		return fmt.Errorf("failed to validate create game: %+v", err)
	}
	db := msg.Context.Get("db").(*mgo.Database)
	game, err := models.CreateGame(db, cg.Name, cg.Type)
	if err != nil {
		models.SendError(msg.Sender.AccountID, userFailMsg)
		return fmt.Errorf("failed to create game: %+v", err)
	}
	if err := models.SendAll(newLobbyMessage(msg.Type, game)); err != nil {
		return fmt.Errorf("sendall error: %+v", err)
	}
	return nil
}

func validateJoinLeaveGame(msg models.Message) (*JoinLeaveGameRequest, error) {
	jlg := struct {
		Game JoinLeaveGameRequest `json:"game"`
	}{}
	err := json.Unmarshal(msg.Raw, &jlg)
	if err != nil {
		return nil, err
	}

	if jlg.Game.ID == "" {
		return nil, fmt.Errorf("game id cannot be empty: %+v", string(msg.Raw))
	}
	return &jlg.Game, nil
}
func HandleJoinGame(msg models.Message) error {
	userFailMsg := "failed to join game"
	jg, err := validateJoinLeaveGame(msg)
	db := msg.Context.Get("db").(*mgo.Database)
	game, err := models.JoinGame(db, jg.ID, msg.Sender)
	if err != nil {
		models.SendError(msg.Sender.AccountID, userFailMsg)
		return fmt.Errorf("%s: %+v ", userFailMsg, err)
	}

	//notify all clients that someone joined this game
	if err := models.SendAll(newLobbyMessage("GAMEJOINALERT", game)); err != nil {
		return fmt.Errorf("sendall error in handleJoinGame: %+v", err)
	}

	//notify this client to enter the game; this ultimately redirs the user into the game
	if err := models.Send(msg.Sender.AccountID, newLobbyMessage("GAMEJOIN", game)); err != nil {
		return fmt.Errorf("send error in handleJoinGame: %+v", err)
	}

	//check to see if game is ready to be started
	/*if err := CheckStartGame(game); err != nil {
		models.SendError(msg.Sender.AccountID, userFailMsg)
		return fmt.Errorf("%s. %+v", userFailMsg, err)
	}*/
	return nil
}

/*
func CheckStartGame(game models.Game) error {
	if game.GameType == "holdem" {
		return CheckStartHoldem(game)
	} else if game.GameType == "highcard" {
		return CheckStartHighCard(game)
	} else {
		return fmt.Errorf("game type: %s, is an invalid gametype", game.GameType)
	}
	return nil
}
*/

func HandleLeaveGame(msg models.Message) error {
	userFailMsg := "failed to leave game"
	lg, err := validateJoinLeaveGame(msg)
	db := msg.Context.Get("db").(*mgo.Database)
	game, err := models.LeaveGame(db, lg.ID, msg.Sender.AccountID)
	if err != nil {
		models.SendError(msg.Sender.AccountID, userFailMsg)
		return fmt.Errorf("%s: %+v", userFailMsg, err.Error())
	}

	//notify all clients that someone left this game
	if err := models.SendAll(newLobbyMessage("GAMELEAVEALERT", game)); err != nil {
		return fmt.Errorf("sendall error in handleLeaveGame: %+v", err)
	}

	//notify this client to leave the game
	if err := models.Send(msg.Sender.AccountID, newLobbyMessage("GAMELEAVE", game)); err != nil {
		return fmt.Errorf("send error in handleLeaveGame: %+v", err)
	}
	return nil
}
