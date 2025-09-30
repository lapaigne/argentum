package main

import (
	"argentum/db"

	"github.com/labstack/echo/v4"
)

func (e Endpoints) menu(c echo.Context) error {
	acc := getClaims(c).Level
	switch acc {
	case db.AccessLevels["worker"]:
		return c.Render(200, "wmenu", nil)
	case db.AccessLevels["dispatcher"]:
		return c.Render(200, "dmenu", nil)
	case db.AccessLevels["admin"]:
		return c.Render(200, "amenu", nil)
	default:
		return echo.ErrUnauthorized
	}
}
