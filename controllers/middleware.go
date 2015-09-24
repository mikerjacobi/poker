package controllers

import (
	"net/http"

	"github.com/Sirupsen/logrus"
	"github.com/labstack/echo"
	"github.com/mikerjacobi/poker/models"
	"github.com/spf13/viper"
	"gopkg.in/mgo.v2"
)

func Auth() echo.HandlerFunc {
	return func(c *echo.Context) error {
		// Skip WebSocket
		if (c.Request().Header.Get(echo.Upgrade)) == echo.WebSocket || c.Request().Method == "OPTIONS" {
			return nil
		}

		sessionID := c.Request().Header.Get("x-session")
		if sessionID == "" {

			//hack for duplicate cookie name...
			authCookie := &http.Cookie{}
			for i := range c.Request().Cookies() {
				cookie := c.Request().Cookies()[i]
				if cookie.Name == "session" && cookie.Value != "" {
					authCookie = cookie
					break
				}
			}

			if authCookie.Name == "" {
				logrus.Errorf("failed to pull session cookie in auth middleware")
				return echo.NewHTTPError(http.StatusUnauthorized)
			}
			sessionID = authCookie.Value
		}

		db := c.Get("db").(*mgo.Database)
		a, err := models.CheckSession(db, sessionID)
		if err != nil {
			logrus.Errorf("failed to check session in auth middleware: %s.", sessionID)
			return echo.NewHTTPError(http.StatusUnauthorized)
		}

		//happy path successful login
		c.Set("user", a)
		return nil
	}
}

func LogStateMiddleware() echo.MiddlewareFunc {
	return func(h echo.HandlerFunc) echo.HandlerFunc {
		return func(c *echo.Context) error {
			logrus.Infof("cookies: %+v", c.Request().Cookies())
			logrus.Infof("headers: %+v", c.Request().Header)
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
			return h(c)
		}
	}
}
