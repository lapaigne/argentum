package main

import (
	"argentum/db"
	"strconv"
	"strings"
	"time"

	"github.com/labstack/echo/v4"
)

func (e Endpoints) me_conf(c echo.Context) error {
	id, err := strconv.Atoi(c.Get("id").(string))
	if err != nil {
		return err
	}

	now := Today()
	if err := db.UpdateTaskStatus(id, SQLTime(now, true), SQLTime(time.Time{}, false)); err != nil {
		return err
	}

	// temp solution, target table or row instead
	uid := getClaims(c).UID
	d := struct {
		Data   *Data
		Helper *Helper
		UID    int
	}{
		Data:   &data,
		Helper: &helper,
		UID:    uid,
	}

	return c.Render(200, "me-tbody", d)
}

func (e Endpoints) me_exp(c echo.Context) error {
	id, err := strconv.Atoi(c.Get("id").(string))
	if err != nil {
		return err
	}

	d := struct {
		Data   *Data
		Helper *Helper
		Id     int
	}{
		Data:   &data,
		Helper: &helper,
		Id:     id,
	}

	return c.Render(200, "mt-exp", d)
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
	uid := getClaims(c).UID

	d := struct {
		Data   *Data
		Helper *Helper
		UID    int
	}{
		Data:   &data,
		Helper: &helper,
		UID:    uid,
	}

	return c.Render(200, "me", d)
}
