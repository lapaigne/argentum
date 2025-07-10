package main

import (
	"argentum/db"
	"fmt"
	"strconv"
	"time"

	"github.com/labstack/echo/v4"
)

func TaskFormHandler(c echo.Context) error {
	var err error
	date := c.FormValue("created")

	var t db.Task

	pdate, err := time.Parse(time.DateOnly, date)
	if err != nil {
		fmt.Println(err)
		return err
	}

	t.Created_date = pdate

	// structurize the categories
	// perhaps add an extra field for level (w/ vals as 1,2 or 3)
	// use proper parser
	t.Cat_1 = 1
	t.Cat_2 = 2
	t.Cat_3 = 3

	t.Desc = c.FormValue("desc")
	// t.Addr_obj
	add := c.FormValue("address")
	temp, err := strconv.Atoi(add)
	if err != nil {
		fmt.Println(err)
		return err
	}
	t.Addr_obj = temp

	t.Comment = c.FormValue("comment")

	fmt.Println(t.Created_date)

	err = db.AddTask(t)

	if err != nil {
		fmt.Println(err)
		fmt.Println(err)
	}

	FetchIncomplete()
	return c.Render(200, "task-list", data.Incomplete)
}

func TaskListHandler(c echo.Context) error {

	FetchIncomplete()
	return c.Render(200, "task-list", data.Incomplete)
}
