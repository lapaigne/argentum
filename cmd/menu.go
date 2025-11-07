package main

import (
	"github.com/labstack/echo/v4"
)

func (e Endpoints) menu(c echo.Context) error {
	acc := getClaims(c).Level
	switch acc {
	case ACC_WORKER:
		return c.Redirect(303, "/me/")
	case ACC_DISPATCHER:
		return c.Render(200, "dmenu", nil)
	case ACC_ADMIN:
		return c.Render(200, "amenu", nil)
	default:
		return echo.ErrUnauthorized
	}
}
