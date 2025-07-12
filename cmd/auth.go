package main

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

func AuthHandler(c echo.Context) error {

	return c.Redirect(http.StatusSeeOther, "/")
}

func AuthTel(c echo.Context) error {
	tel := c.FormValue("tel")
	return c.Render(200, "tel-err", tel)
}
