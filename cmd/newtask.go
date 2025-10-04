package main

import (
	"argentum/db"
	"database/sql"
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/labstack/echo/v4"
)

func (e Endpoints) newtask_GET(c echo.Context) error {

	ctx := struct {
		Data   *Data
		Helper *Helper
		Today  string
		Cat1   sql.NullInt32
		Cat2   sql.NullInt32
		Addr   int
	}{
		Data:   &data,
		Helper: &helper,
		Today:  Today().Format(time.DateOnly),
		Cat1:   SQLInt(-1, false),
		Cat2:   SQLInt(-1, false),
		Addr:   -1,
	}

	return c.Render(200, "newtask", ctx)
}

func (e Endpoints) newtask_POST(c echo.Context) error {
	url := c.Request().URL.Path

	switch url {
	case "/newtask/":
		return e.newtask_(c)
	case "/newtask/act":
		return e.newtask_act(c)
	case "/newtask/addr":
		return e.newtask_addr(c)
	case "/newtask/until":
		return e.newtask_until(c)
	case "/newtask/cat-1":
		return e.newtask_cat1(c)
	case "/newtask/cat-2":
		return e.newtask_cat2(c)
	default:
		fmt.Println(url)
		return echo.ErrNotFound
	}
}

func (e Endpoints) newtask_(c echo.Context) error {

	var err error
	var t db.Task

	date, err := time.Parse(time.DateOnly, c.FormValue("created"))
	if err != nil {
		fmt.Println(err)
		return err
	}

	dif := time.Since(date)

	if dif.Hours() > 24 {
		errS := fmt.Sprintf("Invalid creation date: %v", date)
		return errors.New(errS)
	}

	t.Created_date = date

	c1, err := strconv.Atoi(c.FormValue("cat1"))
	if err != nil {
		return err
	}
	t.Cat1 = c1

	c2, err := strconv.Atoi(c.FormValue("cat2"))
	if err != nil {
		return err
	}
	t.Cat2 = c2

	c3, err := strconv.Atoi(c.FormValue("cat3"))
	if err != nil {
		return err
	}
	t.Cat3 = c3

	t.Desc = c.FormValue("desc")

	addr, err := strconv.Atoi(c.FormValue("address"))
	if err != nil {
		return err
	}
	t.Addr_obj = addr

	t.Comment = c.FormValue("comment")

	w, err := strconv.Atoi(c.FormValue("act"))
	if err != nil {
		return err
	}
	t.Worker = w

	err = db.AddTask(t)
	if err != nil {
		return err
	}

	ctx := struct {
		Data   *Data
		Helper *Helper
	}{
		Data:   &data,
		Helper: &helper,
	}

	if err := data.Fetch(); err != nil {
		return err
	}

	return c.Render(200, "alltasks", ctx)
}

func (e Endpoints) newtask_cat1(c echo.Context) error {

	val, err := strconv.Atoi(c.FormValue("cat1"))
	if err != nil {
		return err
	}

	ctx := struct {
		Data  *Data
		Today string
		Cat1  sql.NullInt32
		Cat2  sql.NullInt32
	}{
		Data:  &data,
		Today: Today().Format(time.DateOnly),
		Cat1:  SQLInt(val, true),
		Cat2:  SQLInt(-1, true),
	}

	return c.Render(200, "tf-cat-1-res", ctx)
}

func (e Endpoints) newtask_cat2(c echo.Context) error {

	val, err := strconv.Atoi(c.FormValue("cat2"))
	if err != nil {
		return err
	}

	ctx := struct {
		Data  *Data
		Today string
		Cat2  sql.NullInt32
	}{
		Data:  &data,
		Today: Today().Format(time.DateOnly),
		Cat2:  SQLInt(val, true),
	}

	return c.Render(200, "tf-cat-3", ctx)
}

func (e Endpoints) newtask_addr(c echo.Context) error {

	val, err := strconv.Atoi(c.FormValue("address"))
	if err != nil {
		return err
	}

	ctx := struct {
		Data   *Data
		Helper *Helper
		Addr   int
	}{
		Data:   &data,
		Helper: &helper,
		Addr:   val,
	}

	return c.Render(200, "addrtable", ctx)
}

func (e Endpoints) newtask_act(c echo.Context) error {
	return c.Render(200, "tf-act", nil)
}

func (e Endpoints) newtask_until(c echo.Context) error {
	return c.Render(200, "tf-until", nil)
}
