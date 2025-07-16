package main

import (
	"argentum/db"
	"fmt"
	"log"
	"time"

	"github.com/golang-jwt/jwt/v5"
	echojwt "github.com/labstack/echo-jwt/v4"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

var data Data

func main() {

	fmt.Println("just so it'd be here")
	e := echo.New()
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

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

	r := e.Group("/d")

	r.Use(JWTCooked())

	config := echojwt.Config{
		NewClaimsFunc: func(c echo.Context) jwt.Claims {
			return new(jwtClaims)
		},
		SigningKey: []byte("secret"),
	}

	r.Use(echojwt.WithConfig(config))

	r.GET("/page", func(c echo.Context) error {
		return c.Render(200, "page", data)
	})

	e.GET("/", func(c echo.Context) error { return c.Render(200, "auth", nil) })
	e.POST("/submit-auth", AuthHandler)

	e.GET("/split", func(c echo.Context) error { return c.Render(200, "split", nil) })
	e.GET("/index", func(c echo.Context) error { return c.Render(200, "index", nil) })

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

	e.GET("/tflist", func(c echo.Context) error {
		return c.Render(200, "tflist", data)
	})

	e.POST("/submit-auth/tel", AuthTel)

	e.Logger.Fatal(e.Start(":42069"))
}
