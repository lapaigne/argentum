package main

import (
	"argentum/db"
	"context"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	echojwt "github.com/labstack/echo-jwt/v4"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

var endpoints Endpoints

var data Data
var helper = Helper{
	NFormat: NFormat,
	Format:  Format,
	Status:  Status,
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
	e.File("/favicon.ico", "assets/img/favicon.ico")

	e.Renderer = templates

	db.OpenConn()
	defer db.CloseConn()

	// init data
	data.fetchUWs(context.Background())
	data.fetchTasks(context.Background())
	data.fetchCats(context.Background())
	data.fetchAddrs(context.Background())

	autoFetch(&data, time.Second*30, time.Second*5)

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

	e.GET("/me/*", endpoints.me_GET)
	e.POST("/me/*", endpoints.me_POST)

	e.GET("/workers/", endpoints.workers_GET)
	e.POST("/workers/*", endpoints.workers_POST)

	e.GET("/newtask/", endpoints.newtask_GET)
	e.POST("/newtask/*", endpoints.newtask_POST)

	e.GET("/cats/", endpoints.cats_GET)
	e.POST("/cats/*", endpoints.cats_POST)

	e.GET("/addrs/", endpoints.addrs_GET)
	e.POST("/addrs/*", endpoints.addrs_POST)

	e.GET("/", endpoints.slash_GET)
	e.POST("/", endpoints.slash_POST)

	e.GET("/menu/", endpoints.menu)

	e.GET("/alltasks/", endpoints.alltasks)

	e.GET("/signin", endpoints.signin)

	e.POST("/signout", endpoints.logout)
	e.POST("/submit-auth", endpoints.login)
	e.POST("/submit-auth/tel", endpoints.authTel)

	e.Logger.Fatal(e.Start(fmt.Sprintf(":%s", os.Getenv("APP_PORT"))))
}
