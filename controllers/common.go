package controllers

import (
	"fmt"
	"text/template"

	"github.com/Sirupsen/logrus"
	"github.com/labstack/echo"
	"github.com/mikerjacobi/poker/models"
	"gopkg.in/mgo.v2"
)

type Response struct {
	Success bool        `json:"success"`
	Payload interface{} `json:"payload"`
}

func HealthCheck(c *echo.Context) error {
	c.String(200, "")
	return nil
}

func Index(c *echo.Context) error {
	dbconn := c.Get("db").(*mgo.Database)
	counts := dbconn.C("counts")

	if err := counts.Insert(&models.TestStruct{"index"}); err != nil {
		c.String(500, fmt.Sprintf("broken: %s", err.Error()))
		return nil
	}

	t, err := template.ParseFiles("static/html/layout.html", "static/html/dashboard.html")
	if err != nil {
		c.String(500, fmt.Sprintf("broken: %s", err.Error()))
		return nil
	}

	user := models.Account{}
	loggedIn := false
	loggedIn, ok := c.Get("logged_in").(bool)
	if ok && loggedIn {
		user = c.Get("user").(models.Account)
	}
	args := map[string]interface{}{
		"Username": user.Username,
		"LoggedIn": loggedIn,
	}
	logrus.Infof("args: %+v", args)
	t.Execute(c.Response(), args)
	return nil
}
