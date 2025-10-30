package main

import (
	"argentum/db"
	"database/sql"
	"fmt"
	"strconv"
	"strings"

	"github.com/labstack/echo/v4"
)

func (e Endpoints) cats_GET(c echo.Context) error {
	if getClaims(c).Level < ACC_DISPATCHER {
		c.Redirect(303, "/menu/")
	}

	return c.Render(200, "cats", data.Cats.Sorted())
}

func (e Endpoints) cats_POST(c echo.Context) error {
	url := c.Request().URL.Path

	if id, ok := strings.CutPrefix(url, "/cats/edit-"); ok {
		c.Set("id", id)
		return e.cats_edit(c)
	}

	if id, ok := strings.CutPrefix(url, "/cats/upd-"); ok {
		c.Set("id", id)
		return e.cats_upd(c)
	}

	if id, ok := strings.CutPrefix(url, "/cats/del-"); ok {
		c.Set("id", id)
		return e.cats_del(c)
	}

	switch url {
	case "/cats/add-panel":
		return e.cats_addpanel(c)
	case "/cats/add":
		return e.cats_add(c)
	default:
		return echo.ErrNotFound
	}
}

func (e Endpoints) cats_addpanel(c echo.Context) error {
	parents := []db.Category{}
	for _, v := range data.Cats {
		if v.Level < 3 {
			parents = append(parents, *v)
		}
	}
	return c.Render(200, "cats-add", parents)
}

func (e Endpoints) cats_edit(c echo.Context) error {
	acc := getClaims(c).Level
	if acc != ACC_ADMIN && acc != ACC_DISPATCHER {
		return echo.ErrUnauthorized
	}

	id, err := strconv.Atoi(c.Get("id").(string))
	if err != nil {
		return err
	}

	d := data.Cats[id]

	return c.Render(200, "cats-upd", d)
}

func (e Endpoints) cats_upd(c echo.Context) error {
	catName := c.FormValue("cat-name")

	id, err := strconv.Atoi(c.Get("id").(string))
	if err != nil {
		return err
	}

	cat := data.Cats[id]
	if catName != cat.Name {
		cat.Name = catName
		if err := db.UpdateCategory(id, *cat); err != nil {
			return err
		}
	}

	return c.Render(200, "cats-row", data.Cats[id])
}

func (e Endpoints) cats_add(c echo.Context) error {
	catName := c.FormValue("cat-name")
	checked := c.FormValue("no-parent")

	par, err := strconv.Atoi(c.FormValue("cats-parent-sel"))
	if err != nil {
		return err
	}

	parCat, ok := data.Cats[par]
	if !ok {
		return fmt.Errorf("no parent cat w/ id %d", par)
	}

	cat := db.Category{
		Name:   catName,
		Active: true,
	}

	var parent sql.NullInt32
	if checked != "checked" {
		parent.Int32 = int32(par)
		parent.Valid = true
		cat.Level = parCat.Level + 1
	} else {
		cat.Level = 1
	}

	cat.Parent = parent

	id, err := db.AddCategory(cat)
	if err != nil {
		return err
	}

	cat.Id = id
	data.Cats[cat.Id] = &cat

	return c.Render(200, "cats-tbody", data.Cats.Sorted())
}

func (e Endpoints) cats_del(c echo.Context) error {
	id, err := strconv.Atoi(c.Get("id").(string))
	if err != nil {
		return err
	}

	cat, ok := data.Cats[id]
	if !ok {
		return echo.ErrNotFound
	}

	cat.Active = !cat.Active

	if err := db.UpdateCategory(id, *cat); err != nil {
		return err
	}

	return c.Render(200, "cats-row", cat)
}
