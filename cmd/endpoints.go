package main

import (
	"argentum/db"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
)

type Endpoints struct{}

var endpoints Endpoints

func (p Endpoints) Setup(e *echo.Echo) error {

	e.GET("/signin", func(c echo.Context) error { return c.Render(200, "signin", nil) })
	e.POST("/", func(c echo.Context) error { return c.Redirect(http.StatusSeeOther, "/") })
	e.GET("/", func(c echo.Context) error { return c.Redirect(http.StatusSeeOther, "/menu") })

	e.POST("/submit-auth", login)

	e.GET("/new-task", func(c echo.Context) error {
		data_old.Helper.Today = time.Now().Format(time.DateOnly)
		return c.Render(200, "tfpage", data_old)
	})

	// e.POST("/submit-task", endpoints.AddTask)

	e.POST("/submit-task", taskFormHandler)
	e.POST("/submit-task/act", ActHandler)
	e.POST("/submit-task/addr", AddrHandler)
	e.POST("/submit-task/until", UntilHandler)
	e.POST("/submit-task/cat-1", Cat1Handler)
	e.POST("/submit-task/cat-2", Cat2Handler)

	e.GET("/task-list", func(c echo.Context) error {
		FetchTasks()
		return c.Render(200, "task-list", data_old)
	})

	e.GET("/auth", func(c echo.Context) error {
		return c.Render(200, "auth", data_old)
	})

	e.POST("/submit-auth/tel", authTel)

	return nil
}

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

func (e Endpoints) AddTask(c echo.Context) error {
	return nil
}

func (e Endpoints) EditAddresses(c echo.Context) error {

	return c.Render(200, "ea-upd", nil)
}
