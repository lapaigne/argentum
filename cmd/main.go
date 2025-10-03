package main

import (
	"argentum/db"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/golang-jwt/jwt/v5"
	echojwt "github.com/labstack/echo-jwt/v4"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

var data_old Data_old
var data Data
var helper = Helper{
	Levels: &db.AccessLevels,
}

var public = map[string]bool{
	"/signin":      true,
	"/submit-auth": true,
}

func main() {
	e := echo.New()
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	e.HTTPErrorHandler = ErrorHandler

	templates, err := NewTemplates("assets")
	if err != nil {
		log.Fatal(err)
	}

	e.Static("/assets", "assets")

	e.Renderer = templates

	db.OpenConn()
	defer db.CloseConn()

	FetchRare()
	FetchTasks()

	if err := data.Fetch(); err != nil {
		fmt.Println(err)
	}

	data_old.Init()
	data_old.Fill()

	config := echojwt.Config{
		SigningKey: []byte(accSecret),
		NewClaimsFunc: func(c echo.Context) jwt.Claims {
			return new(jwtClaims)
		},
		Skipper: func(c echo.Context) bool {
			path := c.Request().URL.Path
			return public[path] || strings.HasPrefix(path, "/assets/")
		},
	}

	e.Use(jwtMiddleware(config))

	e.POST("/signout", logout)

	e.GET("/me/*", endpoints.me_GET)
	e.POST("/me/*", endpoints.me_POST)

	e.GET("/workers/", endpoints.workes_GET)
	e.POST("/workers/*", endpoints.workers_POST)

	e.GET("/newtask/", endpoints.newtask_GET)
	e.POST("/newtask/*", endpoints.newtask_POST)

	e.GET("/menu/", endpoints.menu)

	e.GET("/tflist", func(c echo.Context) error { return c.Render(200, "tflist", data_old) })

	endpoints.Setup(e)

	e.Logger.Fatal(e.Start(fmt.Sprintf(":%s", os.Getenv("APP_PORT"))))
}
