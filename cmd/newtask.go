package main

import (
	"argentum/db"
	"database/sql"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/labstack/echo/v4"
)

func (e Endpoints) newtask_GET(c echo.Context) error {
	if getClaims(c).Level < ACC_DISPATCHER {
		c.Redirect(303, "/menu/")
	}

	d := struct {
		Filled *db.Task
		Errors any
		Data   *Data
		Helper *Helper
		Today  string
		Cat1   sql.NullInt32
		Cat2   sql.NullInt32
		Addr   int
	}{
		Filled: nil,
		Errors: nil,
		Data:   &data,
		Helper: &helper,
		Today:  Today().Format(time.DateOnly),
		Cat1:   SQLInt(-1, false),
		Cat2:   SQLInt(-1, false),
		Addr:   -1,
	}

	return c.Render(200, "newtask", d)
}

func (e Endpoints) newtask_POST(c echo.Context) error {
	url := c.Request().URL.Path

	switch url {
	case "/newtask/":
		return e.newtask_submit(c)
	case "/newtask/addr":
		return e.newtask_addr(c)
	case "/newtask/cat-1":
		return e.newtask_cat1(c)
	case "/newtask/cat-2":
		return e.newtask_cat2(c)
	default:
		fmt.Println(url)
		return echo.ErrNotFound
	}
}

func (e Endpoints) newtask_submit(c echo.Context) error {
	var t db.Task

	errs := struct {
		Any    bool
		Cat1   bool
		Cat2   bool
		Cat3   bool
		Addr   bool
		Worker bool
		Desc   bool
		Until  bool
	}{}

	created, err := time.Parse(time.DateOnly, c.FormValue("created"))
	if err != nil {
		return err
	}

	t.Created_date = created

	var until sql.NullTime

	if c.FormValue("until-checkbox") != "checked" {
		untilTime, err := time.Parse(time.DateOnly, c.FormValue("do-until"))
		if err != nil {
			fmt.Println(err)
			errs.Until = true
			errs.Any = true
			until.Valid = false
		} else {
			until.Time = untilTime
			until.Valid = true
		}
	}

	t.Until_date = until

	c1, err := strconv.Atoi(c.FormValue("cat1"))
	v, ok := data.Cats[c1]
	if err != nil || !ok || v.Level != 1 {
		errs.Cat1 = true
		errs.Cat2 = true
		errs.Cat3 = true
		errs.Any = true
	} else {
		t.Cat1 = c1
	}

	c2, err := strconv.Atoi(c.FormValue("cat2"))
	v, ok = data.Cats[c2]
	if err != nil || !ok || v.Level != 2 {
		errs.Cat2 = true
		errs.Cat3 = true
		errs.Any = true
	} else {
		t.Cat2 = c2
	}

	c3, err := strconv.Atoi(c.FormValue("cat3"))
	v, ok = data.Cats[c2]
	if err != nil || !ok || v.Level != 3 {
		errs.Cat3 = true
		errs.Any = true
	} else {
		t.Cat3 = c3
	}

	desc := c.FormValue("desc")
	if len(strings.TrimSpace(desc)) == 0 {
		errs.Desc = true
		errs.Any = true
	} else {
		t.Desc = desc
	}

	addr, err := strconv.Atoi(c.FormValue("address"))
	if err != nil {
		errs.Addr = true
		errs.Any = true
	} else {
		t.Addr_obj = addr
	}

	t.Comment = c.FormValue("comment")

	w, err := strconv.Atoi(c.FormValue("act"))
	if err != nil {
		errs.Worker = true
		errs.Any = true
	} else {
		t.Worker = w
	}

	d := struct {
		Filled *db.Task
		Errors any
		Data   *Data
		Helper *Helper
		Today  string
		Cat1   sql.NullInt32
		Cat2   sql.NullInt32
		Addr   int
	}{
		Filled: &t,
		Errors: errs,
		Data:   &data,
		Helper: &helper,
		Today:  Today().Format(time.DateOnly),
		Cat1:   SQLInt(t.Cat1, t.Cat1 != 0),
		Cat2:   SQLInt(t.Cat2, t.Cat2 != 0),
		Addr:   t.Addr_obj,
	}

	if errs.Any {
		fmt.Println()
		fmt.Println(&t)
		fmt.Println()
		return c.Render(200, "newtask", d)
	}

	data.MTask++
	if _, err := db.AddTask(t); err != nil {
		fmt.Println(err)
	}
	data.Tasks[data.MTask] = &t

	return c.Redirect(303, "/alltasks/")
}

func (e Endpoints) newtask_cat1(c echo.Context) error {
	val, err := strconv.Atoi(c.FormValue("cat1"))
	if err != nil {
		return err
	}

	d := struct {
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

	return c.Render(200, "tf-cat-1-res", d)
}

func (e Endpoints) newtask_cat2(c echo.Context) error {
	val, err := strconv.Atoi(c.FormValue("cat2"))
	if err != nil {
		return err
	}

	d := struct {
		Data  *Data
		Today string
		Cat2  sql.NullInt32
	}{
		Data:  &data,
		Today: Today().Format(time.DateOnly),
		Cat2:  SQLInt(val, true),
	}

	return c.Render(200, "tf-cat-3", d)
}

func (e Endpoints) newtask_addr(c echo.Context) error {
	val, err := strconv.Atoi(c.FormValue("address"))
	if err != nil {
		return err
	}

	d := struct {
		Data   *Data
		Helper *Helper
		Addr   int
	}{
		Data:   &data,
		Helper: &helper,
		Addr:   val,
	}

	return c.Render(200, "addrtable", d)
}
