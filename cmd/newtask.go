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
	data_old.Helper.Today = time.Now().Format(time.DateOnly)
	return c.Render(200, "newtask", data_old)
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

	if err := data.Fetch(); err != nil {
		fmt.Println(err)
		return err
	}
	data_old.ResetHelper()

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

	c1, err := strconv.Atoi(c.FormValue("cat_1"))
	if err != nil {
		fmt.Println("cat 1")
		return err
	}

	t.Cat_1 = c1

	c2, err := strconv.Atoi(c.FormValue("cat_2"))
	if err != nil {
		fmt.Println("cat 2")
		return err
	}

	t.Cat_2 = c2

	c3, err := strconv.Atoi(c.FormValue("cat_3"))
	if err != nil {
		fmt.Println("cat 3")
		return err
	}

	t.Cat_3 = c3

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
		fmt.Println(err)
		return err
	}

	FetchTasks()

	data_old.Helper.Cat_1 = sql.NullInt32{Int32: int32(0), Valid: false}
	data_old.Helper.Cat_2 = sql.NullInt32{Int32: int32(0), Valid: false}
	data_old.Helper.Addr = -1

	return c.Render(200, "alltasks", data_old)
}

func (e Endpoints) newtask_cat1(c echo.Context) error {

	raw := c.FormValue("cat_1")
	val, err := strconv.Atoi(raw)

	if err != nil {
		return err
	}

	data_old.Helper.Cat_1 = sql.NullInt32{Int32: int32(val), Valid: true}
	data_old.Helper.Cat_2 = sql.NullInt32{Int32: int32(-1), Valid: true}
	return c.Render(200, "tf-cat-1-res", data_old)
}

func (e Endpoints) newtask_cat2(c echo.Context) error {

	raw := c.FormValue("cat_2")
	val, err := strconv.Atoi(raw)

	if err != nil {
		return err
	}

	data_old.Helper.Cat_2 = sql.NullInt32{Int32: int32(val), Valid: true}
	return c.Render(200, "tf-cat-3", data_old)
}

func (e Endpoints) newtask_addr(c echo.Context) error {
	raw := c.FormValue("address")
	val, err := strconv.Atoi(raw)

	if err != nil {
		return err
	}

	data_old.Helper.Addr = val
	return c.Render(200, "addr-table", data_old)
}

func (e Endpoints) newtask_act(c echo.Context) error {
	return c.Render(200, "tf-act", data_old)
}

func (e Endpoints) newtask_until(c echo.Context) error {
	return c.Render(200, "tf-until", data_old)
}
