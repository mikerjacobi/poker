package controllers

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/Sirupsen/logrus"
	"github.com/labstack/echo"
	"github.com/mikerjacobi/poker/server/models"
	"github.com/spf13/viper"
	"gopkg.in/mgo.v2"
)

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

type LobbyController struct {
	DB    *mgo.Database
	Queue chan models.Message
	*models.Comms
	HoldemController
	HighCardController
}

func newLobbyController(db *mgo.Database, comms *models.Comms, hc HoldemController, hcc HighCardController) (LobbyController, error) {
	GameTypes = strings.Split(viper.GetString("game_types"), ",")
	lc := LobbyController{
		DB:                 db,
		Comms:              comms,
		HoldemController:   hc,
		HighCardController: hcc,
	}

	lc.Queue = make(chan models.Message)
	go lc.ReadMessages()
	return lc, nil
}

func (lc LobbyController) ReadMessages() {
	for {
		m := <-lc.Queue
		switch m.Type {
		case models.GameCreate:
			//handle game create
			lc.HandleCreateGame(m)
		//case GameStart:
		//logrus.Infof("game start in lobbyQ readmsgs")
		case models.GameJoin:
			lc.HandleJoinGame(m)
		case models.GameLeave:
			lc.HandleLeaveGame(m)
		default:
			continue
		}
	}
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

	if !models.StringInSlice(cg.Game.Type, GameTypes) {
		return nil, fmt.Errorf("invalid gametype: %s", cg.Game.Type)
	}
	return &cg.Game, nil
}

func (lc LobbyController) HandleCreateGame(msg models.Message) {
	log := logrus.WithFields(logrus.Fields{"func": "HandleCreateGame"})
	cg, err := validateCreateGame(msg)
	if err != nil {
		e := "failed to validate create game "
		sendError(lc.Comms, msg.SenderAccountID, e)
		logrus.Errorf("%s: %s", msg.SenderAccountID, e+err.Error())
		return
	}
	game, err := models.CreateGame(lc.DB, cg.Name, cg.Type)
	if err != nil {
		e := "failed to create game "
		sendError(lc.Comms, msg.SenderAccountID, e)
		logrus.Errorf("%s: %s", msg.SenderAccountID, e+err.Error())
		return
	}
	resp := models.LobbyMessage{Message: msg, Game: game}
	if err := lc.SendAll(resp); err != nil {
		log.Errorf("sendall error: %+v", err)
		return
	}
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
func (lc LobbyController) HandleJoinGame(msg models.Message) {
	log := logrus.WithFields(logrus.Fields{"func": "HandleJoinGame"})
	jg, err := validateJoinLeaveGame(msg)

	account, ok := msg.Context.Get("user").(models.Account)
	if !ok {
		e := "failed to get account from context in handleJoinGame"
		sendError(lc.Comms, msg.SenderAccountID, e)
		logrus.Errorf("%s: %s", msg.SenderAccountID, e)
		return
	}
	game, err := models.JoinGame(lc.DB, jg.ID, account)
	if err != nil {
		e := "failed to join game. "
		sendError(lc.Comms, msg.SenderAccountID, e)
		logrus.Errorf("%s: %s", msg.SenderAccountID, e+err.Error())
		return
	}

	//notify all clients that someone joined this game
	resp := models.LobbyMessage{Message: models.Message{Type: models.GameJoinAlert}, Game: game}
	if err := lc.SendAll(resp); err != nil {
		log.Errorf("sendall error: %+v", err)
		return
	}

	//notify this client to enter the game
	resp = models.LobbyMessage{Message: models.Message{Type: models.GameJoin}, Game: game}
	if err := lc.Send(msg.SenderAccountID, resp); err != nil {
		log.Errorf("sendall error: %+v", err)
		return
	}

	//check to see if game is ready to be started
	if err := lc.CheckStartGame(game); err != nil {
		log.Warnf("game not started, continuing. %+v", err)
	}
}

func (lc LobbyController) CheckStartGame(game models.Game) error {
	if game.GameType == "holdem" {
		return lc.HoldemController.CheckStartGame(game)
	} else if game.GameType == "highcard" {
		return lc.HighCardController.CheckStartGame(game)
	} else {
		return fmt.Errorf("game type: %s, is an invalid gametype", game.GameType)
	}
	return nil
}

func (lc LobbyController) HandleLeaveGame(msg models.Message) {
	log := logrus.WithFields(logrus.Fields{"func": "HandleLeaveGame"})
	lg, err := validateJoinLeaveGame(msg)
	game, err := models.LeaveGame(lc.DB, lg.ID, msg.SenderAccountID)
	if err != nil {
		e := "failed to leave game. "
		sendError(lc.Comms, msg.SenderAccountID, e)
		logrus.Errorf("%s: %s", msg.SenderAccountID, e+err.Error())
		return
	}

	//notify all clients that someone left this game
	resp := models.LobbyMessage{Message: models.Message{Type: models.GameLeaveAlert}, Game: game}
	if err := lc.SendAll(resp); err != nil {
		log.Errorf("sendall error: %+v", err)
		return
	}

	//notify this client to leave the game
	resp = models.LobbyMessage{Message: models.Message{Type: models.GameLeave}, Game: game}
	if err := lc.Send(msg.SenderAccountID, resp); err != nil {
		log.Errorf("sendall error: %+v", err)
		return
	}
}
