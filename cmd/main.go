package main

import (
	"argentum/db"
	"fmt"
	"log"

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

	e.HTTPErrorHandler = ErrorHandler

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

	config := echojwt.Config{
		NewClaimsFunc: func(c echo.Context) jwt.Claims {
			return new(jwtClaims)
		},
		SigningKey: []byte("secret"),
	}

	w := e.Group("/w")

	w.Use(JWTCooked())
	w.Use(AutoRefreshJWT)
	w.Use(echojwt.WithConfig(config))
	w.Use(JWTRoles(db.AccessLevels["worker"]))

	w.GET("/menu", func(c echo.Context) error { return c.Render(200, "wmenu", nil) })

	d := e.Group("/d")

	d.Use(JWTCooked())
	d.Use(AutoRefreshJWT)
	d.Use(echojwt.WithConfig(config))
	d.Use(JWTRoles(db.AccessLevels["dispatcher"]))

	d.GET("/menu", func(c echo.Context) error { return c.Render(200, "dmenu", nil) })

	d.GET("/tflist", func(c echo.Context) error {
		return c.Render(200, "tflist", data)
	})

	a := e.Group("/a")

	a.Use(JWTCooked())
	a.Use(AutoRefreshJWT)
	a.Use(echojwt.WithConfig(config))
	a.Use(JWTRoles(db.AccessLevels["admin"]))

	a.GET("/menu", func(c echo.Context) error { return c.Render(200, "amenu", nil) })

	endpoints.Setup(e)

	e.Logger.Fatal(e.Start(":42069"))
}
