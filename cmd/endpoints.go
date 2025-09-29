package main

import (
	"argentum/db"
	"database/sql"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/labstack/echo/v4"
)

type Endpoints struct{}

var endpoints Endpoints

func (p Endpoints) Setup(e *echo.Echo) error {

	e.GET("/signin", func(c echo.Context) error { return c.Render(200, "signin", nil) })
	e.POST("/", func(c echo.Context) error { return c.Redirect(http.StatusSeeOther, "/") })
	e.GET("/", func(c echo.Context) error { return c.Redirect(http.StatusSeeOther, "/menu") })

	e.POST("/submit-auth", login)

	e.GET("/new-task", func(c echo.Context) error {
		data_old.Helper.Today = time.Now().Format(time.DateOnly)
		return c.Render(200, "tfpage", data_old)
	})

	// e.POST("/submit-task", endpoints.AddTask)

	e.POST("/submit-task", taskFormHandler)
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

	e.POST("/submit-auth/tel", authTel)

	return nil
}

func (p Endpoints) menu(c echo.Context) error {
	acc := getClaims(c).Level
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

	acc := getClaims(c).Level
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
		Helper Helper0
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

	acc := getClaims(c).Level
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

func (p Endpoints) editWorkers(c echo.Context) error {

	acc := getClaims(c).Level
	if acc != db.AccessLevels["admin"] {
		return echo.ErrUnauthorized
	}

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return err
	}

	uw := data.UWs[id]
	d := struct {
		W      db.Worker
		U      db.User
		Target string
	}{
		W:      uw.W,
		U:      uw.U,
		Target: fmt.Sprintf("ew-%d", id),
	}

	return c.Render(200, "ew-upd", d)
}

func (p Endpoints) AddTask(c echo.Context) error {
	return nil
}

func (p Endpoints) EditAddresses(c echo.Context) error {

	return c.Render(200, "ea-upd", nil)
}
