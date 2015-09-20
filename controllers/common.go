package controllers

import (
	"fmt"
	"text/template"

	"github.com/Sirupsen/logrus"
	"github.com/labstack/echo"
	"github.com/mikerjacobi/echomongo/models"
	"github.com/spf13/viper"
	"gopkg.in/mgo.v2"
)

type Response struct {
	Success bool `json:"success"`
	Payload interface{}
}

func HealthCheck(c *echo.Context) error {
	c.String(200, "")
	return nil
}

func Index(c *echo.Context) error {

	authCookie, err := c.Request().Cookie("testcook")
	logrus.Infof(">>> cooki: %+v, err: %+v", authCookie, err)

	dbconn := c.Get("db").(*mgo.Database)
	counts := dbconn.C("counts")

	if err := counts.Insert(&models.TestStruct{"index"}); err != nil {
		c.String(500, fmt.Sprintf("broken: %s", err.Error()))
		return nil
	}

	t, err := template.ParseFiles("static/html/layout.html", "static/html/greet.html", "static/html/mainPage.html")
	if err != nil {
		c.String(500, fmt.Sprintf("broken: %s", err.Error()))
		return nil
	}

	loggedIn := false
	user, ok := c.Get("user").(models.Account)
	if ok {
		loggedIn = user.Username != ""
	}
	args := map[string]interface{}{
		"Username": user.Username,
		"LoggedIn": loggedIn,
		"Logout":   fmt.Sprintf("http://username:password@%s", viper.GetString("base_uri"))}
	t.Execute(c.Response(), args)
	return nil
}
