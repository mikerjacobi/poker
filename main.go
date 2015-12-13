package main

import (
	"flag"
	"net/http"

	"github.com/Sirupsen/logrus"
	"github.com/labstack/echo"
	mw "github.com/labstack/echo/middleware"
	"github.com/mikerjacobi/poker/controllers"
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
	viper.BindEnv("POKER_MONGO_1_PORT_27017_TCP_ADDR")

	//setup config
	viper.AddConfigPath(configpath)
	viper.AddConfigPath("/")
	viper.SetConfigName(config)
	viper.ReadInConfig()
}

func main() {
	// database
	session, err := mgo.Dial(viper.GetString("POKER_MONGO_1_PORT_27017_TCP_ADDR"))
	if err != nil {
		panic(err)
	}
	defer session.Close()
	db := session.Clone().DB(viper.GetString("database"))
	if err := controllers.InitializeMessageHandler(db); err != nil {
		logrus.Panicf("failed to init message handler")
	}

	router := echo.New()

	// middleware
	router.Use(controllers.LogStateMiddleware())
	router.Use(mw.Logger())
	router.Use(mw.Recover())
	router.Use(controllers.CORSMiddleware())
	router.Use(controllers.MongoMiddleware(session))
	router.Use(controllers.CheckAuth())

	// unauthed routes
	router.Index("static/index.html")
	router.Favicon("static/favicon.ico")
	router.Get("/geta", controllers.GetA)
	router.Get("/getb", controllers.GetB)
	router.Post("/create_account", controllers.CreateAccount)
	router.Get("/hello", hello)
	router.Post("/login", controllers.Login)
	router.Static("/s/", "static")

	// auth routes
	auth := router.Group("")
	auth.Use(controllers.RequireAuth())
	auth.Get("/authhello", authhello)
	auth.Get("/math", controllers.GetMathCount)
	auth.WebSocket("/ws", controllers.HandleWebSocket)
	auth.Post("/logout", controllers.Logout)
	auth.Post("/game", controllers.CreateGame)
	auth.Get("/game/:gameID", controllers.GetGame)
	auth.Get("/games", controllers.GetOpenGames)
	auth.Post("/game/:gameID/join", controllers.JoinGame)

	// start server
	logrus.Info("starting server")
	router.Run("0.0.0.0:80")
}
