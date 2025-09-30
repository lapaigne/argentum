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

func (e Endpoints) me_(c echo.Context) error {
	data_old.Fill()

	claims := getClaims(c)
	uid := claims.UID

	ctx := struct {
		Data   Data
		Helper Helper
	}{
		Data:   data,
		Helper: helper,
	}

	fmt.Println(ctx)

	ctx0 := struct {
		Raw    RawData0
		Mapped MappedData0
		Helper Helper0
		UID    int
	}{
		Raw:    data_old.Raw,
		Mapped: data_old.Mapped,
		Helper: data_old.Helper,
		UID:    uid,
	}

	return c.Render(200, "me", ctx0)
}

func (e Endpoints) me_conf(c echo.Context) error {
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
		return err
	}

	if ok {
		if err := FetchTasks(); err != nil {
			fmt.Println(err)
		}

		if err := data_old.Fill(); err != nil {
			fmt.Println(err)
		}
	}

	// return c.Render(200, "me", nil)
	return c.NoContent(http.StatusNoContent)
}

func (e Endpoints) me_exp(c echo.Context) error {
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

func (e Endpoints) me_POST(c echo.Context) error {
	url := c.Request().URL.Path
	switch url {
	case "/me/conf-:id":
		return e.me_conf(c)
	case "/me/exp-:id":
		return e.me_exp(c)
	default:
		return echo.ErrNotFound
	}
}

func (e Endpoints) me_GET(c echo.Context) error {
	url := c.Request().URL.Path
	switch url {
	case "/me/":
		return e.me_(c)
	default:
		return echo.ErrNotFound
	}
}
