package main

import (
	"argentum/db"
	"argentum/emb"
	"fmt"
	"log"
	"strings"

	"github.com/golang-jwt/jwt/v5"
	echojwt "github.com/labstack/echo-jwt/v4"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

var data Data

var public = map[string]bool{
	"/signin":      true,
	"/submit-auth": true,
	"/":            true,
}

func main() {

	fmt.Println("just so it'd be here")
	e := echo.New()
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	e.HTTPErrorHandler = ErrorHandler

	e.Static("/css", "views/css")

	templates, err := emb.NewTemplates("views")
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
		Skipper: func(c echo.Context) bool {
			path := c.Request().URL.Path
			return public[path] || strings.HasPrefix(path, "/css/")
		},
	}

	e.Use(
		JWTCooked(),
		AutoRefreshJWT,
		echojwt.WithConfig(config),
	)

	e.GET("/me", func(c echo.Context) error { return c.Render(200, "me", data) })

	e.POST("/me/conf-:id", endpoints.Me_ConfirmTask)
	e.POST("/me/exp-:id", endpoints.Me_ExpandTask)

	e.GET("/menu", endpoints.Menu)

	e.GET("/tflist", func(c echo.Context) error { return c.Render(200, "tflist", data) })

	endpoints.Setup(e)

	e.Logger.Fatal(e.Start(":42069"))
}
