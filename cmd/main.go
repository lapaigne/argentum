package main

import (
	"argentum/db"
	"fmt"
	"log"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

var data Data

func main() {

	fmt.Println("just so it'd be here")
	e := echo.New()
	e.Use(middleware.Logger())

	e.Static("/css", "views/css")

	templates, err := NewTemplates("views")
	if err != nil {
		log.Fatal(err)
	}

	e.Renderer = templates

	db.OpenConn()
	defer db.CloseConn()

	FetchRare()
	FetchTasks()

	data.Init()
	data.Fill()

	e.GET("/split", func(c echo.Context) error { return c.Render(200, "split", nil) })
	e.GET("/", func(c echo.Context) error { return c.Render(200, "index", nil) })

	e.GET("/new-task", func(c echo.Context) error {
		data.Helper.Today = time.Now().Format(time.DateOnly)
		if c.Request().Header.Get("HX-Request") == "true" {
			return c.Render(200, "task-form", data)
		}
		return c.Render(200, "new-task", data)
	})

	e.GET("/table", func(c echo.Context) error {
		if c.Request().Header.Get("HX-Request") == "true" {
			return c.Render(200, "table", data)
		}
		return c.Render(200, "base", data)
	})

	e.GET("/page2", func(c echo.Context) error {
		return c.Render(200, "page2", data)
	})

	e.GET("/page", func(c echo.Context) error {
		return c.Render(200, "page", data)
	})

	e.POST("/submit-task", TaskFormHandler)
	e.POST("/submit-task/act", ActHandler)
	e.POST("/submit-task/cat-1", Cat1Handler)
	e.POST("/submit-task/cat-2", Cat2Handler)

	e.GET("/task-list", func(c echo.Context) error {
		FetchTasks()
		return c.Render(200, "task-list", data)
	})

	e.GET("/auth", func(c echo.Context) error {
		return c.Render(200, "auth", data)
	})

	e.POST("/submit-auth", AuthHandler)
	e.POST("/submit-auth/tel", AuthTel)

	e.Logger.Fatal(e.Start(":42069"))
}
