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
	"/mes/":        true,
}

func main() {
	e := echo.New()
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	e.HTTPErrorHandler = ErrorHandler

	templates, err := NewTemplates("views")
	if err != nil {
		log.Fatal(err)
	}

	e.Static("/css", "views/css")

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
			return public[path] || strings.HasPrefix(path, "/css/")
		},
	}

	e.Use(jwtMiddleware(config))

	e.POST("/signout", logout)

	e.GET("/me", endpoints.me)

	e.POST("/me/conf-:id", endpoints.Me_ConfirmTask)
	e.POST("/me/exp-:id", endpoints.Me_ExpandTask)

	e.GET("/menu", endpoints.menu)

	e.GET("/tflist", func(c echo.Context) error { return c.Render(200, "tflist", data_old) })

	e.GET("/edit-workers", func(c echo.Context) error { return c.Render(200, "edit-workers", data) })

	e.GET("/mes/*", endpoints.me_s)

	e.POST("/edit-workers/edit-:id", endpoints.editWorkers)

	e.POST("/edit-workers/upd", func(c echo.Context) error {
		f := c.FormValue("f_name")
		i := c.FormValue("i_name")
		o := c.FormValue("o_name")

		d := UserWorker{
			W: db.Worker{
				F_name: f,
				I_name: i,
				O_name: o,
			},
			U: db.User{},
		}

		return c.Render(200, "ew-row", d)
	})

	e.POST("/edit-workers/add-panel", func(c echo.Context) error {

		if err := isAdminErr(c); err != nil {
			return err
		}

		return c.Render(200, "ew-add", nil)
	})

	endpoints.Setup(e)

	e.Logger.Fatal(e.Start(fmt.Sprintf(":%s", os.Getenv("APP_PORT"))))
}
