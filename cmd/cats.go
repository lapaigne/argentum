package main

import (
	"strings"

	"github.com/labstack/echo/v4"
)

func (e Endpoints) cats_GET(c echo.Context) error {
	return c.Render(200, "cats", data.Cats)
}

func (e Endpoints) cats_POST(c echo.Context) error {
	url := c.Request().URL.Path

	if id, ok := strings.CutPrefix(url, "/cats/edit-"); ok {
		c.Set("id", id)
		return e.cats_edit(c)
	}

	if id, ok := strings.CutPrefix(url, "/cats/upd-"); ok {
		c.Set("id", id)
		return e.cats_upd(c)
	}

	if id, ok := strings.CutPrefix(url, "/cats/del-"); ok {
		c.Set("id", id)
		return e.cats_del(c)
	}

	switch url {
	case "/cats/add-panel":
		return e.cats_addpanel(c)
	case "/cats/add":
		return e.cats_add(c)
	default:
		return echo.ErrNotFound
	}
}

// click on btn, calls panel
func (e Endpoints) cats_addpanel(c echo.Context) error {
	return nil
}

// click on row, calls panel
func (e Endpoints) cats_edit(c echo.Context) error {
	return nil
}

// returns row (upd data)
func (e Endpoints) cats_upd(c echo.Context) error {
	return nil
}

// returns row (new row)
func (e Endpoints) cats_add(c echo.Context) error {
	return nil
}

// returns row (soft del)
func (e Endpoints) cats_del(c echo.Context) error {
	return nil
}
