package main

import (
	"argentum/db"
	"database/sql"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
)

type Endpoints struct{}

var endpoints Endpoints

func (p Endpoints) Setup(e *echo.Echo) error {

	e.GET("/signin", func(c echo.Context) error { return c.Render(200, "signin", nil) })
	e.POST("/", func(c echo.Context) error { return c.Redirect(http.StatusSeeOther, "/") })

	e.GET("/", func(c echo.Context) error {
		cookie, err := c.Cookie("token")
		if err != nil || cookie == nil {
			return c.Redirect(http.StatusSeeOther, "/signin")
		}

		token, err := jwt.ParseWithClaims(cookie.Value, &jwtClaims{}, func(t *jwt.Token) (any, error) {
			return []byte("secret"), nil
		})

		if err != nil || !token.Valid {
			return c.Redirect(http.StatusSeeOther, "/signin")
		}
		return c.Redirect(http.StatusSeeOther, "/menu")
	})

	e.POST("/submit-auth", Login)

	e.GET("/new-task", func(c echo.Context) error {
		data_old.Helper.Today = time.Now().Format(time.DateOnly)
		return c.Render(200, "tfpage", data_old)
	})

	// e.POST("/submit-task", endpoints.AddTask)

	e.POST("/submit-task", TaskFormHandler)
	e.POST("/submit-task/act", ActHandler)
	e.POST("/submit-task/addr", AddrHandler)
	e.POST("/submit-task/until", UntilHandler)
	e.POST("/submit-task/cat-1", Cat1Handler)
	e.POST("/submit-task/cat-2", Cat2Handler)

	e.GET("/task-list", func(c echo.Context) error {
		FetchTasks()
		return c.Render(200, "task-list", data_old)
	})

	e.GET("/auth", func(c echo.Context) error {
		return c.Render(200, "auth", data_old)
	})

	e.POST("/submit-auth/tel", AuthTel)

	return nil
}

func (p Endpoints) Menu(c echo.Context) error {
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
}

func (p Endpoints) Me_ExpandTask(c echo.Context) error {

	acc := GetClaims(c).Level
	switch acc {
	case db.AccessLevels["worker"]:
	case db.AccessLevels["dispatcher"]:
	case db.AccessLevels["admin"]:
	default:
		return echo.ErrUnauthorized
	}

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return err
	}

	FetchTasks()

	ctx := struct {
		UID    int
		Task   *db.Task
		Cats   map[int]string
		Addrs  map[int]string
		Helper Helper
	}{
		UID:    id,
		Task:   data_old.Mapped.Tasks[id],
		Cats:   data_old.Mapped.Categories,
		Addrs:  data_old.Mapped.Addresses,
		Helper: data_old.Helper,
	}

	return c.Render(200, "mt-exp", ctx)
}

func (p Endpoints) Me_ConfirmTask(c echo.Context) error {

	acc := GetClaims(c).Level
	switch acc {
	case db.AccessLevels["worker"]:
	case db.AccessLevels["dispatcher"]:
	case db.AccessLevels["admin"]:
	default:
		return echo.ErrUnauthorized
	}

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return err
	}

	now := time.Now().Round(time.Hour * 24)
	ok, err := db.UpdateTask(id, sql.NullTime{Time: now, Valid: true}, sql.NullTime{Valid: false})

	if err != nil {
		println(err)
	}

	if ok {

		if err := FetchTasks(); err != nil {
			fmt.Println(err)
		}

		if err := data_old.Fill(); err != nil {
			fmt.Println(err)
		}

		fmt.Println("updated")
	}

	return c.NoContent(http.StatusNoContent)
	// return c.Render(200, "me", nil)
}

func (p Endpoints) EditWorkers(c echo.Context) error {

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return err
	}

	w := data_old.Mapped.ProperWorkers[id]
	d := struct {
		Worker *db.Worker
		Target string
	}{
		Worker: w,
		Target: fmt.Sprintf("ew-%d", id),
	}

	acc := GetClaims(c).Level
	if acc != db.AccessLevels["admin"] {
		return echo.ErrUnauthorized
	}

	return c.Render(200, "ew-upd", d)
}

func (p Endpoints) AddTask(c echo.Context) error {
	return nil
}
