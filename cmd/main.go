package main

import (
	"argentum/db"
	"database/sql"
	"fmt"
	"html/template"
	"io"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

type Templates struct {
	templates *template.Template
}

func (t *Templates) Render(w io.Writer, name string, data any, c echo.Context) error {
	return t.templates.ExecuteTemplate(w, name, data)
}

func NewTemplate() *Templates {
	return &Templates{
		templates: template.Must(template.ParseGlob("views/*.html")),
	}
}

var data Data

func main() {

	data.Helper.Cat_1 = sql.NullInt32{Int32: -1, Valid: true}
	data.Helper.Cat_2 = sql.NullInt32{Int32: -1, Valid: true}

	fmt.Println("just so it'd be here")
	e := echo.New()
	e.Use(middleware.Logger())

	e.Static("/css", "views/css")

	e.Renderer = NewTemplate()

	db.OpenConn()
	defer db.CloseConn()

	FetchRare()
	FetchTasks()

	data.Init()
	data.Fill()

	e.Router()

	e.GET("/", func(c echo.Context) error { return c.Render(200, "index", nil) })

	e.GET("/task-form", func(c echo.Context) error {
		data.Helper.Today = time.Now().Format(time.DateOnly)
		return c.Render(200, "task-form", data)
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
