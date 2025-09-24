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

func authTel(c echo.Context) error {
	tel := c.FormValue("tel")
	return c.Render(200, "tel-err", tel)
}

func taskFormHandler(c echo.Context) error {

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

	return c.Render(200, "tflist", data_old)
}

func Cat1Handler(c echo.Context) error {

	raw := c.FormValue("cat_1")
	val, err := strconv.Atoi(raw)

	if err != nil {
		return err
	}

	data_old.Helper.Cat_1 = sql.NullInt32{Int32: int32(val), Valid: true}
	data_old.Helper.Cat_2 = sql.NullInt32{Int32: int32(-1), Valid: true}
	return c.Render(200, "tf-cat-1-res", data_old)
}

func Cat2Handler(c echo.Context) error {

	raw := c.FormValue("cat_2")
	val, err := strconv.Atoi(raw)

	if err != nil {
		return err
	}

	data_old.Helper.Cat_2 = sql.NullInt32{Int32: int32(val), Valid: true}
	return c.Render(200, "tf-cat-3", data_old)
}

func AddrHandler(c echo.Context) error {
	raw := c.FormValue("address")
	val, err := strconv.Atoi(raw)

	if err != nil {
		return err
	}

	data_old.Helper.Addr = val
	return c.Render(200, "addr-table", data_old)
}

func ActHandler(c echo.Context) error {
	return c.Render(200, "tf-act", data_old)
}

func UntilHandler(c echo.Context) error {
	return c.Render(200, "tf-until", data_old)
}
