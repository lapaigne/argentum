package main

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

type Endpoints struct{}

var endpoints Endpoints

func (p Endpoints) Setup(e *echo.Echo) error {

	e.GET("/signin", func(c echo.Context) error { return c.Render(200, "signin", nil) })
	e.POST("/", func(c echo.Context) error { return c.Redirect(http.StatusSeeOther, "/") })
	e.GET("/", func(c echo.Context) error { return c.Redirect(http.StatusSeeOther, "/menu/") })

	e.POST("/submit-auth", login)

	e.GET("/task-list", func(c echo.Context) error {
		FetchTasks()
		return c.Render(200, "task-list", data_old)
	})

	e.GET("/auth", func(c echo.Context) error {
		return c.Render(200, "auth", data_old)
	})

	e.POST("/submit-auth/tel", p.authTel)

	return nil
}

func (e Endpoints) AddTask(c echo.Context) error {
	return nil
}

func (e Endpoints) EditAddresses(c echo.Context) error {

	return c.Render(200, "ea-upd", nil)
}
