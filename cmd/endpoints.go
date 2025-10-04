package main

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

type Endpoints struct{}

func (e Endpoints) slash_GET(c echo.Context) error {
	return c.Redirect(http.StatusSeeOther, "/menu/")
}

func (e Endpoints) slash_POST(c echo.Context) error {
	return c.Redirect(http.StatusSeeOther, "/menu/")
}
