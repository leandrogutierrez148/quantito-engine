package controllers

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

type PingController struct {
}

func NewPingController() *PingController {
	return &PingController{}
}

func (controller *PingController) Ping(c echo.Context) error {
	return c.String(http.StatusOK, "Pong")
}
