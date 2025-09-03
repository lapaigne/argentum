package main

import (
	"sort"

	"github.com/labstack/echo/v4"
)

type Sorting int8

const (
	None Sorting = iota
	Ascending
	Descending
)

type TaskDisplay struct {
	Adrr_obj       Sorting
	Created_date   Sorting
	Until_date     Sorting
	Mark_date      Sorting
	Completed_date Sorting
	Worker         Sorting
}

func SortHandler(c echo.Context) error {
	sk := c.QueryParam("sort")
	// order := c.QueryParam("order")

	tasks := data_old.Mapped.Tasks

	sort.Slice(tasks, func(i, j int) bool {
		var less bool
		switch sk {
		case "cat":
			less = tasks[i].Created_date.Before(tasks[j].Created_date)
		case "addr":
			less = tasks[i].Created_date.Before(tasks[j].Created_date)
		case "created":
			less = tasks[i].Created_date.Before(tasks[j].Created_date)
		case "until":
			less = tasks[i].Until_date.Time.Before(tasks[j].Until_date.Time)
		case "mark":
			less = tasks[i].Mark_date.Time.Before(tasks[j].Mark_date.Time)
		case "completed":
			less = tasks[i].Completed_date.Time.Before(tasks[j].Completed_date.Time)
		default:
			less = tasks[i].Created_date.Before(tasks[j].Created_date)
		}

		return less
	})
	return nil
}
