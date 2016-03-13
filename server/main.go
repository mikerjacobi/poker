package main

import (
	"flag"
	"net/http"

	"github.com/Sirupsen/logrus"
	"github.com/labstack/echo"
	mw "github.com/labstack/echo/middleware"
	"github.com/mikerjacobi/poker/server/controllers"
	"github.com/spf13/viper"
	"gopkg.in/mgo.v2"
)

// Handler
func hello(c *echo.Context) error {
	return c.String(http.StatusOK, "unauth2 Hello, World!\n")
}
func authhello(c *echo.Context) error {
	h := c.Request().Header.Get("Authorization")
	logrus.Infof("header: %+v", h)
	return c.String(http.StatusOK, "auth Hello, World!\n")
}

var (
	config     string
	configpath string
)

func init() {
	flag.StringVar(&config, "config", "config", "config name [.toml,.json,.yml]")
	flag.StringVar(&configpath, "configpath", ".", "config location")
	flag.Parse()

	//env
	viper.BindEnv("SERVER_MONGO_1_PORT_27017_TCP_ADDR")

	//setup config
	viper.AddConfigPath(configpath)
	viper.AddConfigPath("/")
	viper.SetConfigName(config)
	viper.ReadInConfig()
}

func main() {
	// database
	session, err := mgo.Dial(viper.GetString("SERVER_MONGO_1_PORT_27017_TCP_ADDR"))
	if err != nil {
		panic(err)
	}
	defer session.Close()
	db := session.Clone().DB(viper.GetString("database"))

	router := echo.New()

	// middleware
	router.Use(controllers.LogStateMiddleware())
	router.Use(mw.Logger())
	router.Use(mw.Recover())
	router.Use(controllers.CORSMiddleware())
	router.Use(controllers.MongoMiddleware(session))
	router.Use(controllers.CheckAuth())

	// unauthed routes
	router.Index("../client/index.html")
	router.Favicon("../client/favicon.ico")
	router.Get("/geta", controllers.GetA)
	router.Get("/getb", controllers.GetB)
	router.Post("/account", controllers.CreateAccount)
	router.Get("/hello", hello)
	router.Post("/login", controllers.Login)
	router.Static("/", "../client")

	// auth routes
	auth := router.Group("")
	auth.Use(controllers.RequireAuth())
	auth.Get("/authhello", authhello)
	auth.Get("/math", controllers.GetMathCount)
	auth.Post("/logout", controllers.Logout)
	auth.Get("/games", controllers.GetOpenGames)
	auth.Get("/game/:gameID", controllers.GetGame)

	//websocket actions
	mh, err := controllers.InitializeMessageHandler(db)
	if err != nil {
		logrus.Panicf("failed to init message handler")
	}
	auth.WebSocket("/ws", mh.HandleWebSocket)
	mh.Handle("/default", controllers.DefaultActionHandler)
	mh.Handle("/ws/connect", controllers.HandleWebSocketConnect)
	mh.Handle("/ws/disconnect", controllers.HandleWebSocketDisconnect)

	mh.Handle("/math/increment", controllers.HandleIncrement)
	mh.Handle("/math/decrement", controllers.HandleDecrement)
	mh.Handle("/math/square", controllers.HandleSquare)
	mh.Handle("/math/sqrt", controllers.HandleSqrt)

	mh.Handle("/game/create", controllers.HandleCreateGame)
	mh.Handle("/game/join", controllers.HandleJoinGame)
	mh.Handle("/game/leave", controllers.HandleLeaveGame)

	mh.Handle("/highcard/play", controllers.HandlePlay)

	mh.Handle("/account/load", controllers.HandleLoadAccount)
	mh.Handle("/account/chips/request", controllers.HandleChipRequest)
	//start server
	logrus.Info("starting server")
	router.Run("0.0.0.0:80")
}
