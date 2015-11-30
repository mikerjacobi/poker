package controllers

import (
	"net/http"

	"github.com/Sirupsen/logrus"
	"github.com/labstack/echo"
	"github.com/mikerjacobi/poker/models"
	"github.com/spf13/viper"
	"gopkg.in/mgo.v2"
)

func CheckAuth() echo.HandlerFunc {
	return func(c *echo.Context) error {
		loggedIn := false
		sessionID := ""

		//if this is a websocket connection, check cookies
		if (c.Request().Header.Get(echo.Upgrade)) == echo.WebSocket {
			sessionCookie, err := c.Request().Cookie("session")
			if err != nil {
				return err
			}
			sessionID = sessionCookie.Value
		} else {
			sessionID = c.Request().Header.Get("x-session")
		}

		if db, ok := c.Get("db").(*mgo.Database); sessionID != "" && ok {
			account, err := models.CheckSession(db, sessionID)
			if err != nil {
				logrus.Errorf("failed to check session in auth middleware: %s.", sessionID)
				return nil
			}

			//happy path successful login
			loggedIn = true
			c.Set("user", account)
		}
		c.Set("logged_in", loggedIn)
		return nil
	}
}
func RequireAuth() echo.HandlerFunc {
	return func(c *echo.Context) error {
		loggedIn, ok := c.Get("logged_in").(bool)
		if ok && loggedIn {
			return nil
		}
		return echo.NewHTTPError(http.StatusUnauthorized)
	}
}

func LogStateMiddleware() echo.MiddlewareFunc {
	return func(h echo.HandlerFunc) echo.HandlerFunc {
		return func(c *echo.Context) error {
			return h(c)
		}
	}
}

func MongoMiddleware(session *mgo.Session) echo.MiddlewareFunc {
	return func(h echo.HandlerFunc) echo.HandlerFunc {
		return func(c *echo.Context) error {
			clone := session.Clone()
			defer clone.Close()
			c.Set("db", clone.DB(viper.GetString("database")))
			return h(c)
		}
	}
}

func CORSMiddleware() echo.MiddlewareFunc {
	return func(h echo.HandlerFunc) echo.HandlerFunc {
		return func(c *echo.Context) error {
			c.Response().Header().Set("Access-Control-Allow-Origin", "*")
			c.Response().Header().Set("Access-Control-Allow-Methods", "*")
			c.Response().Header().Set("Access-Control-Allow-Headers", "x-session")
			if c.Request().Method == "OPTIONS" {
				return HealthCheck(c)
			}
			return h(c)
		}
	}
}
