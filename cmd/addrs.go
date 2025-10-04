package main

import "github.com/labstack/echo/v4"

func (e Endpoints) EditAddresses(c echo.Context) error {
	return c.Render(200, "ea-upd", nil)
}
