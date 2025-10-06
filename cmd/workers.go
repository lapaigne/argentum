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

	if id, ok := strings.CutPrefix(url, "/workers/upd-"); ok {
		c.Set("id", id)
		return e.workers_upd(c)
	}

	switch url {
	case "/workers/add-panel":
		return e.workers_addpanel(c)
	case "/workers/add":
		return e.workers_add(c)
	case "/workers/del":
		return e.workers_del(c)
	default:
		return echo.ErrNotFound
	}
}

func (e Endpoints) workers_add(c echo.Context) error {
	return nil
}

func (e Endpoints) workers_del(c echo.Context) error {
	return nil
}

func (e Endpoints) workers_upd(c echo.Context) error {
	f := c.FormValue("f_name")
	i := c.FormValue("i_name")
	o := c.FormValue("o_name")
	l := c.FormValue("ew-lvl-sel")

	wid, err := strconv.Atoi(c.Get("id").(string))
	if err != nil {
		return err
	}

	uw := data.UWs[wid]
	uid := uw.U.Id

	level, err := strconv.Atoi(l)
	if err != nil {
		return err
	}

	d := UserWorker{
		W: db.Worker{
			F_name: f,
			I_name: i,
			O_name: o},
		U: db.User{
			Level: level,
		},
	}


	// slow for now, perhaps would be updated to use channels
	dw, du := Diffs(*uw, d)
	if dw {
		if err = db.UpdateWorker(wid, d.W); err != nil {
			fmt.Println(err)
		}
	}
	if du {
		if err = db.UpdateUserLevel(uid, d.U.Level); err != nil {
			fmt.Println(err)
		}
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

	type Opt struct {
		Level int
		Name  string
		Sel   bool
	}

	uw := data.UWs[id]

	var opts []Opt

	switch uw.U.Level {

	case 0:
		opts = []Opt{
			{Level: 0, Name: "Отключенный аккаунт", Sel: true},
			{Level: 10, Name: "Работник", Sel: false},
			{Level: 50, Name: "Диспетчер", Sel: false},
			{Level: 100, Name: "Админ", Sel: false},
		}
	case 10:
		opts = []Opt{
			{Level: 10, Name: "Работник", Sel: true},
			{Level: 50, Name: "Диспетчер", Sel: false},
			{Level: 100, Name: "Админ", Sel: false},
		}
	case 50:
		opts = []Opt{
			{Level: 10, Name: "Работник", Sel: false},
			{Level: 50, Name: "Диспетчер", Sel: true},
			{Level: 100, Name: "Админ", Sel: false},
		}
	case 100:
		opts = []Opt{
			{Level: 10, Name: "Работник", Sel: false},
			{Level: 50, Name: "Диспетчер", Sel: false},
			{Level: 100, Name: "Админ", Sel: true},
		}
	}

	d := struct {
		W    db.Worker
		U    db.User
		Id   int
		Opts []Opt
	}{
		W:    uw.W,
		U:    uw.U,
		Id:   id,
		Opts: opts,
	}

	return c.Render(200, "ew-upd", d)
}

func (e Endpoints) workers_addpanel(c echo.Context) error {
	if err := isAdminErr(c); err != nil {
		return err
	}
	return c.Render(200, "ew-add", nil)
}
