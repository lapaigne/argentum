package main

import (
	"argentum/db"
	"fmt"
	"log"
	"net/http"

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

	e.GET("/me", func(c echo.Context) error {
		return c.Render(200, "me", data)
	},
		JWTCooked(),
		AutoRefreshJWT,
		echojwt.WithConfig(config),
	)

	e.POST("/me/:id", endpoints.W_ConfirmTask)

	e.GET("/menu", func(c echo.Context) error {
		acc := GetClaims(c).Level
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

	},
		JWTCooked(),
		AutoRefreshJWT,
		echojwt.WithConfig(config),
	)

	e.GET("/tflist", func(c echo.Context) error {
		return c.Render(200, "tflist", data)
	})

	endpoints.Setup(e)

	e.Logger.Fatal(e.Start(":42069"))
}
