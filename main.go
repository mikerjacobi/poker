package main

import (
	"flag"
	"fmt"
	"net/http"

	log "github.com/Sirupsen/logrus"
	"github.com/labstack/echo"
	mw "github.com/labstack/echo/middleware"
	"github.com/mikerjacobi/echomongo/controllers"
	"github.com/spf13/viper"
	"gopkg.in/mgo.v2"
)

// Handler
func hello(c *echo.Context) error {
	return c.String(http.StatusOK, "unauth2 Hello, World!\n")
}
func authhello(c *echo.Context) error {
	h := c.Request().Header.Get("Authorization")
	log.Infof("header: %+v", h)
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
	viper.BindEnv("ECHOMONGO_MONGO_1_PORT_27017_TCP_ADDR")

	//setup config
	viper.AddConfigPath(configpath)
	viper.AddConfigPath("/")
	viper.SetConfigName(config)
	viper.ReadInConfig()
}

func main() {
	// Echo instance
	e := echo.New()

	// Middleware
	e.Use(controllers.LogStateMiddleware())
	e.Use(mw.Logger())
	e.Use(mw.Recover())
	e.Use(controllers.CORSMiddleware())

	mongoIP := viper.GetString("ECHOMONGO_MONGO_1_PORT_27017_TCP_ADDR")
	session, err := mgo.Dial(mongoIP)
	if err != nil {
		panic(err)
	}
	defer session.Close()
	e.Use(controllers.MongoMiddleware(session))

	// unauthed routes
	e.Get("/", controllers.Index)
	e.Post("/create_account", controllers.CreateAccount)
	e.Get("/hello", hello)
	e.Post("/login", controllers.Login)
	e.Get("/favicon.ico", hello)

	// auth routes
	a := e.Group("/a")
	a.Use(controllers.Auth())
	a.Get("/authhello", authhello)
	a.Post("/logout", controllers.Logout)
	a.Options("/logout", controllers.HealthCheck)
	a.Get("/", controllers.Index)

	// static
	e.Static("/s/", "static")

	// Start server
	fmt.Println()
	fmt.Println()
	fmt.Println()
	log.Info("starting server")
	e.Run("0.0.0.0:80")
}
