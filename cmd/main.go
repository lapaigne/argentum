package main

import (
	"argentum/db"
	"argentum/emb"
	"fmt"
	"log"
	"os"
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

	e.POST("/signout", Logout)

	e.GET("/css/:fn", func(c echo.Context) error {
		fn := c.Param("fn")
		f, err := emb.EFS.ReadFile("views/css/" + fn)
		if err != nil {
			fmt.Println(err)
			return echo.ErrNotFound
		}

		return c.Blob(200, "text/css", f)
	})

	e.GET("/me", func(c echo.Context) error {

		data.Fill()

		user := c.Get("user").(*jwt.Token)
		claims := user.Claims.(*jwtClaims)
		uid := claims.UID

		ctx := struct {
			Raw    RawData
			Mapped MappedData
			Helper Helper
			UID    int
		}{
			Raw:    data.Raw,
			Mapped: data.Mapped,
			Helper: data.Helper,
			UID:    uid,
		}

		return c.Render(200, "me", ctx)
	})

	e.POST("/me/conf-:id", endpoints.Me_ConfirmTask)
	e.POST("/me/exp-:id", endpoints.Me_ExpandTask)

	e.GET("/menu", endpoints.Menu)

	e.GET("/tflist", func(c echo.Context) error { return c.Render(200, "tflist", data) })

	e.GET("/edit-workers", func(c echo.Context) error { return c.Render(200, "edit-workers", data) })

	e.POST("/edit-workers/edit-:id", endpoints.EditWorkers)

	e.POST("/edit-workers/upd", func(c echo.Context) error {
		f := c.FormValue("f_name")
		i := c.FormValue("i_name")
		o := c.FormValue("o_name")

		d := db.Worker{
			F_name: f,
			I_name: i,
			O_name: o,
		}

		return c.Render(200, "ew-row", d)
	})

	endpoints.Setup(e)

	e.Logger.Fatal(e.Start(fmt.Sprintf(":%s", os.Getenv("APP_PORT"))))

}
