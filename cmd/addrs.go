package main

import (
	"argentum/db"
	"strconv"
	"strings"

	"github.com/labstack/echo/v4"
)

func (e Endpoints) addrs_GET(c echo.Context) error {
	if getClaims(c).Level < ACC_DISPATCHER {
		c.Redirect(303, "/menu/")
	}

	return c.Render(200, "addrs", data.Addrs)
}

func (e Endpoints) addrs_POST(c echo.Context) error {
	url := c.Request().URL.Path

	if id, ok := strings.CutPrefix(url, "/addrs/edit-"); ok {
		c.Set("id", id)
		return e.addrs_edit(c)
	}

	if id, ok := strings.CutPrefix(url, "/addrs/upd-"); ok {
		c.Set("id", id)
		return e.addrs_upd(c)
	}

	if id, ok := strings.CutPrefix(url, "/addrs/del-"); ok {
		c.Set("id", id)
		return e.addrs_del(c)
	}

	switch url {
	case "/addrs/add-panel":
		return e.addrs_addpanel(c)
	case "/addrs/add":
		return e.addrs_add(c)
	default:
		return echo.ErrNotFound
	}
}

func (e Endpoints) addrs_addpanel(c echo.Context) error {
	return c.Render(200, "addrs-add", nil)
}

func (e Endpoints) addrs_edit(c echo.Context) error {
	acc := getClaims(c).Level
	if acc != ACC_ADMIN && acc != ACC_DISPATCHER {
		return echo.ErrUnauthorized
	}

	id, err := strconv.Atoi(c.Get("id").(string))
	if err != nil {
		return err
	}

	d := data.Addrs[id]

	return c.Render(200, "addrs-upd", d)
}

func (e Endpoints) addrs_upd(c echo.Context) error {
	addr := c.FormValue("addr")

	id, err := strconv.Atoi(c.Get("id").(string))
	if err != nil {
		return err
	}

	a := data.Addrs[id]
	if addr != a.Address {
		a.Address = addr
		if err := db.UpdateAddress(id, *a); err != nil {
			return err
		}
	}

	return c.Render(200, "addrs-row", data.Addrs[id])
}

func (e Endpoints) addrs_add(c echo.Context) error {
	addr := c.FormValue("addr")

	a := db.Address{
		Address: addr,
		Active:  true,
	}

	id, err := db.AddAddress(a)
	if err != nil {
		return err
	}

	a.Id = id
	data.Addrs[a.Id] = &a

	return c.Render(200, "addrs-row", a)
}

func (e Endpoints) addrs_del(c echo.Context) error {
	id, err := strconv.Atoi(c.Get("id").(string))
	if err != nil {
		return err
	}

	a, ok := data.Addrs[id]
	if !ok {
		return echo.ErrNotFound
	}

	a.Active = false

	if err := db.UpdateAddress(id, *a); err != nil {
		return err
	}

	return c.Render(200, "addrs-row", a)
}
