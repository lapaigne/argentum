package main

import (
	"argentum/db"
	"database/sql"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/labstack/echo/v4"
)

func TaskFormHandler(c echo.Context) error {
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

	// structurize the categories
	// perhaps add an extra field for level (w/ vals as 1,2 or 3)
	// use proper parser

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

	err = db.AddTask(t)

	if err != nil {
		fmt.Println(err)
		return err
	}

	FetchIncomplete()

	return c.Redirect(http.StatusSeeOther, "/task-list")

	// return c.Render(200, "task-list", data.Incomplete)
}

func Cat1Handler(c echo.Context) error {

	raw := c.FormValue("cat_1")
	val, err := strconv.Atoi(raw)

	if err != nil {
		return err
	}

	data.Cat_1 = sql.NullInt32{Int32: int32(val), Valid: true}
	data.Cat_2 = sql.NullInt32{Int32: int32(-1), Valid: true}
	return c.Render(200, "tf-cat-1-res", data)
}

func Cat2Handler(c echo.Context) error {

	raw := c.FormValue("cat_2")
	val, err := strconv.Atoi(raw)

	if err != nil {
		return err
	}

	fmt.Println(data.Cat_2.Int32)
	data.Cat_2 = sql.NullInt32{Int32: int32(val), Valid: true}
	return c.Render(200, "tf-cat-3", data)
}
