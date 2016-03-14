package controllers

import (
	"encoding/json"
	"fmt"

	"github.com/Sirupsen/logrus"
	"github.com/labstack/echo"
	"github.com/mikerjacobi/poker/server/models"
	"gopkg.in/mgo.v2"
)

const (
	gameJoin       = "/game/join"
	gameJoinAlert  = "/game/join/alert"
	gameLeave      = "/game/leave"
	gameLeaveAlert = "/game/leave/alert"
)

type LobbyMessage struct {
	Type        string `json:"type"`
	models.Game `json:"game"`
}

type GameMessage struct {
	Type   string `json:"type"`
	GameID string `json:"gameID"`
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
type GameRequest struct {
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
	games, err := models.LoadGames(db)
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

type JoinGameRequest struct {
	Game        GameRequest `json:"game"`
	BuyinAmount int         `json:buyinAmount`
}

func validateJoinGame(msg models.Message) (*JoinGameRequest, error) {
	jg := JoinGameRequest{}
	err := json.Unmarshal(msg.Raw, &jg)
	if err != nil {
		return nil, err
	}

	if jg.Game.ID == "" {
		return nil, fmt.Errorf("game id cannot be empty: %+v", string(msg.Raw))
	}
	return &jg, nil
}

func validateLeaveGame(msg models.Message) (*GameRequest, error) {
	lg := struct {
		Game GameRequest `json:"game"`
	}{}
	err := json.Unmarshal(msg.Raw, &lg)
	if err != nil {
		return nil, err
	}

	if lg.Game.ID == "" {
		return nil, fmt.Errorf("game id cannot be empty: %+v", string(msg.Raw))
	}
	return &lg.Game, nil
}
func HandleJoinGame(msg models.Message) error {
	userFailMsg := "failed to join game"
	jg, err := validateJoinGame(msg)
	if err != nil {
		models.SendError(msg.Sender.AccountID, userFailMsg)
		return fmt.Errorf("fail validate join game.  %s: %+v ", userFailMsg, err)
	}

	db := msg.Context.Get("db").(*mgo.Database)
	game, err := models.JoinGame(db, jg.Game.ID, msg.Sender.AccountID, jg.BuyinAmount)
	if err != nil {
		models.SendError(msg.Sender.AccountID, userFailMsg)
		return fmt.Errorf("fail join game.  %s: %+v ", userFailMsg, err)
	}

	//notify all clients that someone joined this game
	if err := models.SendAll(newLobbyMessage(gameJoinAlert, game)); err != nil {
		return fmt.Errorf("sendall error in handleJoinGame: %+v", err)
	}

	//notify this client to enter the game; this ultimately redirs the user into the game
	if err := models.Send(msg.Sender.AccountID, newLobbyMessage(gameJoin, game)); err != nil {
		return fmt.Errorf("send error in handleJoinGame: %+v", err)
	}

	return nil
}

func HandleLeaveGame(msg models.Message) error {
	userFailMsg := "failed to leave game"
	lg, err := validateLeaveGame(msg)
	if err != nil {
		models.SendError(msg.Sender.AccountID, userFailMsg)
		return fmt.Errorf("validate %s: %+v", userFailMsg, err.Error())
	}

	db := msg.Context.Get("db").(*mgo.Database)
	game, err := models.LeaveGame(db, lg.ID, msg.Sender.AccountID)
	if err != nil {
		models.SendError(msg.Sender.AccountID, userFailMsg)
		return fmt.Errorf("leavegame %s: %+v", userFailMsg, err.Error())
	}

	//notify all clients that someone left this game
	if err := models.SendAll(newLobbyMessage(gameLeaveAlert, game)); err != nil {
		return fmt.Errorf("sendall error in handleLeaveGame: %+v", err)
	}

	//notify this client to leave the game
	if err := models.Send(msg.Sender.AccountID, newLobbyMessage(gameLeave, game)); err != nil {
		return fmt.Errorf("send error in handleLeaveGame: %+v", err)
	}
	return nil
}

func HandleGameBuyIn(msg models.Message) error {
	return nil
}
func HandleGameCashOut(msg models.Message) error {
	return nil
}
