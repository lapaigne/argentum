package main

import (
	"argentum/db"
	"fmt"
	"html/template"
	"io"

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

	fmt.Println("just so it'd be here")
	e := echo.New()
	e.Use(middleware.Logger())

	e.Static("/css", "views/css")

	e.Renderer = NewTemplate()

	db.OpenConn()
	defer db.CloseConn()

	FetchRare()
	FetchIncomplete()

	fmt.Println(len(data.Incomplete))
	fmt.Println(len(data.Workers))
	fmt.Println(len(data.Categories))
	fmt.Println(len(data.Addresses))

	e.Router()

	e.GET("/", func(c echo.Context) error { return c.Render(200, "index", nil) })
	e.GET("/input", func(c echo.Context) error { return c.Render(200, "input", data) })

	e.POST("/add", func(c echo.Context) error { return c.Redirect(303, "input") })
	e.POST("/task", TaskFormHandler)

	e.POST("/task-list", TaskListHandler)
	e.GET("/task-list", TaskListHandler)

	e.Logger.Fatal(e.Start("0.0.0.0:42069"))
}
