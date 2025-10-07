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

	if id, ok := strings.CutPrefix(url, "/addrs/upd-"); ok {
		c.Set("id", id)
		return e.addrs_upd(c)
	}

	if id, ok := strings.CutPrefix(url, "/addrs/del-"); ok {
		c.Set("id", id)
		return e.addrs_del(c)
	}

	switch url {
	case "/addrs/add-panel":
		return e.addrs_addpanel(c)
	case "/addrs/add":
		return e.addrs_add(c)
	default:
		return echo.ErrNotFound
	}
}

// click on btn, calls panel
func (e Endpoints) addrs_addpanel(c echo.Context) error {
	return nil
}

// click on row, calls panel
func (e Endpoints) addrs_edit(c echo.Context) error {
	return nil
}

// returns row (upd data)
func (e Endpoints) addrs_upd(c echo.Context) error {
	return nil
}

// returns row (new row)
func (e Endpoints) addrs_add(c echo.Context) error {
	return nil
}

// returns row (soft del)
func (e Endpoints) addrs_del(c echo.Context) error {
	return nil
}
