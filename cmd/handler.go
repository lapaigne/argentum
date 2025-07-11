package main

import (
	"argentum/db"
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/labstack/echo/v4"
)

func TaskFormHandler(c echo.Context) error {
	var err error
	var t db.Task

	date := c.FormValue("created")

	pdate, err := time.Parse(time.DateOnly, date)
	if err != nil {
		fmt.Println(err)
		return err
	}

	dif := time.Since(pdate)

	if dif.Hours() > 24 {
		errS := fmt.Sprintf("Invalid creation date: %v", pdate)
		return errors.New(errS)
	}

	t.Created_date = pdate

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
	return c.Render(200, "task-list", data.Incomplete)
}

func Cat1Handler(c echo.Context) error {
	raw := c.FormValue("cat_1")
	cat_1, err := strconv.Atoi(raw)

	// put into html {{ $cat_1 }} instead

	// perhaps replace only the data which would go there

	fmt.Println(cat_1)

	if err != nil {
		return err
	}

	return c.Render(200, "cat-2", cat_1)
}
