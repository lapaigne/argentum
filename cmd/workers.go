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
		return e.workersEdit(c)
	}

	switch url {
	case "/workers/upd":
		return e.workersUpd(c)
	default:
		fmt.Println(url)
		return echo.ErrNotFound
	}
}

func (e Endpoints) workersUpd(c echo.Context) error {
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

func (e Endpoints) workersEdit(c echo.Context) error {

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
