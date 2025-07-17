package main

import (
	"time"

	"github.com/labstack/echo/v4"
)

type Endpoints struct{}

var endpoints Endpoints

func (p Endpoints) Setup(e *echo.Echo) error {

	e.GET("/", func(c echo.Context) error { return c.Render(200, "auth", nil) })
	e.POST("/submit-auth", AuthHandler)

	e.GET("/menu", func(c echo.Context) error { return c.Render(200, "dmenu", nil) })

	e.GET("/new-task", func(c echo.Context) error {
		data.Helper.Today = time.Now().Format(time.DateOnly)
		return c.Render(200, "tfpage", data)
	})

	e.POST("/submit-task", endpoints.AddTask)

	e.POST("/submit-task", TaskFormHandler)
	e.POST("/submit-task/act", ActHandler)
	e.POST("/submit-task/addr", AddrHandler)
	e.POST("/submit-task/until", UntilHandler)
	e.POST("/submit-task/cat-1", Cat1Handler)
	e.POST("/submit-task/cat-2", Cat2Handler)

	e.GET("/task-list", func(c echo.Context) error {
		FetchTasks()
		return c.Render(200, "task-list", data)
	})

	e.GET("/auth", func(c echo.Context) error {
		return c.Render(200, "auth", data)
	})

	e.POST("/submit-auth/tel", AuthTel)

	return nil
}

func (p Endpoints) AddTask(c echo.Context) error {
	return nil
}
