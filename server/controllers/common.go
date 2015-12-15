package controllers

import "github.com/labstack/echo"

type Response struct {
	Success bool        `json:"success"`
	Payload interface{} `json:"payload"`
}

func HealthCheck(c *echo.Context) error {
	c.String(200, "")
	return nil
}

type GetData struct {
	Data string `json:"data"`
}

func GetA(c *echo.Context) error {
	g := GetData{"aAaAAaaAaaaAaA"}
	return c.JSON(200, Response{true, g})
}
func GetB(c *echo.Context) error {
	g := GetData{"bbBBBbBBbBbbBBbb"}
	return c.JSON(200, Response{true, g})
}
