package main

import (
	"argentum/db"
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

	e.GET("/task-form", func(c echo.Context) error {
		data.Today = time.Now().Format(time.DateOnly)
		return c.Render(200, "input", data)
	})

	e.POST("/task-form-submit", TaskFormHandler)

	e.GET("/task-list", func(c echo.Context) error {
		FetchIncomplete()
		return c.Render(200, "task-list", data.Incomplete)
	})

	e.Logger.Fatal(e.Start(":42069"))
}
