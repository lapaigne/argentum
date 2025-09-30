package main

import (
	"argentum/db"
	"fmt"
	"strconv"
	"strings"

	"github.com/labstack/echo/v4"
)

func (e Endpoints) workes_GET(c echo.Context) error {
	return c.Render(200, "workers", data)
}

func (e Endpoints) workers_POST(c echo.Context) error {
	url := c.Request().URL.Path

	if id, ok := strings.CutPrefix(url, "/workers/edit-"); ok {
		c.Set("id", id)
		return e.workers_edit(c)
	}

	switch url {
	case "/workers/add-panel":
		return e.workers_add(c)
	case "/workers/upd":
		return e.workers_upd(c)
	default:
		fmt.Println(url)
		return echo.ErrNotFound
	}
}

func (e Endpoints) workers_upd(c echo.Context) error {
	f := c.FormValue("f_name")
	i := c.FormValue("i_name")
	o := c.FormValue("o_name")

	d := UserWorker{
		W: db.Worker{
			F_name: f,
			I_name: i,
			O_name: o},
		U: db.User{},
	}

	return c.Render(200, "ew-row", d)
}

func (e Endpoints) workers_edit(c echo.Context) error {

	acc := getClaims(c).Level
	if acc != db.AccessLevels["admin"] && acc != db.AccessLevels["dispatcher"] {
		return echo.ErrUnauthorized
	}

	id, err := strconv.Atoi(c.Get("id").(string))
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

func (e Endpoints) workers_add(c echo.Context) error {
	if err := isAdminErr(c); err != nil {
		return err
	}
	return c.Render(200, "ew-add", nil)
}
