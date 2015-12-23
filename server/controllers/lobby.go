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
	//lobby actions
	GameCreate     = "GAMECREATE"
	GameStart      = "GAMESTART"
	GameJoin       = "GAMEJOIN"
	GameJoinAlert  = "GAMEJOINALERT"
	GameLeave      = "GAMELEAVE"
	GameLeaveAlert = "GAMELEAVEALERT"
	LobbyActions   = []string{GameCreate, GameStart, GameJoin, GameJoinAlert, GameLeave, GameLeaveAlert}
	GameTypes      []string
)

type CreateGameRequest struct {
	Name string `json:"gameName"`
	Type string `json:"gameType"`
}
type JoinLeaveGameRequest struct {
	ID string `json:"gameID"`
}
type LobbyMessage struct {
	Message
	models.Game `json:"game"`
}

type LobbyController struct {
	DB    *mgo.Database
	Queue chan Message
	*models.Comms
}

func newLobbyController(db *mgo.Database, comms *models.Comms) (LobbyController, error) {
	GameTypes = strings.Split(viper.GetString("game_types"), ",")
	lc := LobbyController{
		DB:    db,
		Comms: comms,
	}

	lc.Queue = make(chan Message)
	go lc.ReadMessages()
	return lc, nil
}

func (lc LobbyController) ReadMessages() {
	for {
		m := <-lc.Queue
		switch m.Type {
		case GameCreate:
			//handle game create
			lc.HandleCreateGame(m)
		//case GameStart:
		//logrus.Infof("game start in lobbyQ readmsgs")
		case GameJoin:
			lc.HandleJoinGame(m)
		case GameLeave:
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

func validateCreateGame(msg Message) (*CreateGameRequest, error) {
	cg := struct {
		Game CreateGameRequest `json:"game"`
	}{}
	err := json.Unmarshal(msg.Raw, &cg)
	if err != nil {
		return nil, err
	}
	logrus.Errorf("%+v", cg)

	if cg.Game.Name == "" {
		return nil, fmt.Errorf("gamename cannot be empty: %+v", string(msg.Raw))
	}

	if !models.StringInSlice(cg.Game.Type, GameTypes) {
		return nil, fmt.Errorf("invalid gametype: %s", cg.Game.Type)
	}
	return &cg.Game, nil
}

func (lc LobbyController) HandleCreateGame(msg Message) {
	log := logrus.WithFields(logrus.Fields{"func": "HandleCreateGame"})
	cg, err := validateCreateGame(msg)
	if err != nil {
		e := "failed to validate create game "
		sendError(lc.Comms, msg.WebSocketID, e)
		logrus.Errorf("%s: %s", msg.Sender.AccountID, e+err.Error())
		return
	}
	game, err := models.CreateGame(lc.DB, cg.Name, cg.Type)
	if err != nil {
		e := "failed to create game "
		sendError(lc.Comms, msg.WebSocketID, e)
		logrus.Errorf("%s: %s", msg.Sender.AccountID, e+err.Error())
		return
	}
	resp := LobbyMessage{Message: msg, Game: game}
	if err := lc.SendAll(resp); err != nil {
		log.Errorf("sendall error: %+v", err)
		return
	}
}

func validateJoinLeaveGame(msg Message) (*JoinLeaveGameRequest, error) {
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
func (lc LobbyController) HandleJoinGame(msg Message) {
	log := logrus.WithFields(logrus.Fields{"func": "HandleJoinGame"})
	jg, err := validateJoinLeaveGame(msg)
	game, err := models.JoinGame(lc.DB, jg.ID, msg.Sender)
	if err != nil {
		e := "failed to join game. "
		sendError(lc.Comms, msg.WebSocketID, e)
		logrus.Errorf("%s: %s", msg.Sender.AccountID, e+err.Error())
		return
	}

	//notify all clients that someone joined this game
	resp := LobbyMessage{Message: Message{Type: GameJoinAlert}, Game: game}
	if err := lc.SendAll(resp); err != nil {
		log.Errorf("sendall error: %+v", err)
		return
	}

	//notify this client to enter the game
	resp = LobbyMessage{Message: Message{Type: GameJoin}, Game: game}
	if err := lc.Send(msg.WebSocketID, resp); err != nil {
		log.Errorf("sendall error: %+v", err)
		return
	}
}

func (lc LobbyController) HandleLeaveGame(msg Message) {
	log := logrus.WithFields(logrus.Fields{"func": "HandleLeaveGame"})
	lg, err := validateJoinLeaveGame(msg)
	game, err := models.LeaveGame(lc.DB, lg.ID, msg.Sender.AccountID)
	if err != nil {
		e := "failed to leave game. "
		sendError(lc.Comms, msg.WebSocketID, e)
		logrus.Errorf("%s: %s", msg.Sender.AccountID, e+err.Error())
		return
	}

	//notify all clients that someone left this game
	resp := LobbyMessage{Message: Message{Type: GameLeaveAlert}, Game: game}
	if err := lc.SendAll(resp); err != nil {
		log.Errorf("sendall error: %+v", err)
		return
	}

	//notify this client to leave the game
	resp = LobbyMessage{Message: Message{Type: GameLeave}, Game: game}
	if err := lc.Send(msg.WebSocketID, resp); err != nil {
		log.Errorf("sendall error: %+v", err)
		return
	}
}
