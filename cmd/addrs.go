package main

import (
	"strings"

	"github.com/labstack/echo/v4"
)

func (e Endpoints) addrs_GET(c echo.Context) error {
	return c.Render(200, "addrs", data)
}

func (e Endpoints) addrs_POST(c echo.Context) error {
	url := c.Request().URL.Path

	if id, ok := strings.CutPrefix(url, "/addrs/edit-"); ok {
		c.Set("id", id)
		return e.addrs_edit(c)
	}

	switch url {
	case "/addrs/add-panel":
		return e.addrs_addpanel(c)
	case "/addrs/upd":
		return e.addrs_upd(c)
	default:
		return echo.ErrNotFound
	}

}

func (e Endpoints) addrs_addpanel(c echo.Context) error {
	return nil
}

func (e Endpoints) addrs_edit(c echo.Context) error {
	return nil
}

func (e Endpoints) addrs_upd(c echo.Context) error {
	return c.Render(200, "addrs-upd", nil)
}

func (e Endpoints) addrs_del(c echo.Context) error {
	return nil
}
