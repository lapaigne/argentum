package main

import (
	"github.com/labstack/echo/v4"
)

type Endpoints struct{}

func (e Endpoints) slash_GET(c echo.Context) error {
	return c.Redirect(303, "/menu/")
}

func (e Endpoints) slash_POST(c echo.Context) error {
	return c.Redirect(303, "/menu/")
}
