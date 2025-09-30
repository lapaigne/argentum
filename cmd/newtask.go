package main

import (
	"fmt"
	"time"

	"github.com/labstack/echo/v4"
)

func (e Endpoints) newtask_GET(c echo.Context) error {
	data_old.Helper.Today = time.Now().Format(time.DateOnly)
	return c.Render(200, "newtask", data_old)
}

func (e Endpoints) newtask_POST(c echo.Context) error {
	url := c.Request().URL.Path

	switch url {
	case "/newtask/":
		return e.taskFormHandler(c)
	case "/newtask/act":
		return e.actHandler(c)
	case "/newtask/addr":
		return e.addrHandler(c)
	case "/newtask/until":
		return e.untilHandler(c)
	case "/newtask/cat-1":
		return e.cat1Handler(c)
	case "/newtask/cat-2":
		return e.cat2Handler(c)
	default:
		fmt.Println(url)
		return echo.ErrNotFound
	}
}
