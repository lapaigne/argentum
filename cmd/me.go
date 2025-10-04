package main

import (
	"argentum/db"
	"database/sql"
	"fmt"
	"net/http"
	"strconv"
	"strings"
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
		UID    int
	}{
		Data:   data,
		Helper: helper,
		UID:    uid,
	}

	return c.Render(200, "me", ctx)
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

	id, err := strconv.Atoi(c.Get("id").(string))
	if err != nil {
		return err
	}

	now := time.Now().Truncate(time.Hour * 24)
	ok, err := db.UpdateTask(id, sql.NullTime{Time: now, Valid: true}, sql.NullTime{Valid: false})

	if err != nil {
		fmt.Println(err)
		return err
	}

	if ok {
		if err := FetchTasks(); err != nil {
			fmt.Println(err)
			return err
		}

		if err := data_old.Fill(); err != nil {
			fmt.Println(err)
			return err
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

	id, err := strconv.Atoi(c.Get("id").(string))
	if err != nil {
		return err
	}

	FetchTasks()

	ctx := struct {
		Data   *Data
		Helper *Helper
		Id     int
	}{
		Data:   &data,
		Helper: &helper,
		Id:     id,
	}

	return c.Render(200, "mt-exp", ctx)
}

func (e Endpoints) me_POST(c echo.Context) error {
	url := c.Request().URL.Path

	if id, ok := strings.CutPrefix(url, "/me/conf-"); ok {
		c.Set("id", id)
		return e.me_conf(c)
	}

	if id, ok := strings.CutPrefix(url, "/me/exp-"); ok {
		c.Set("id", id)
		return e.me_exp(c)
	}

	switch url {
	default:
		return echo.ErrNotFound
	}
}

func (e Endpoints) me_GET(c echo.Context) error {
	return e.me_(c)
}
