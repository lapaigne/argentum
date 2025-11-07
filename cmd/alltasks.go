package main

import (
	"argentum/db"
	"strconv"
	"time"

	"github.com/labstack/echo/v4"
)

type fStatus int

const (
	FILTER_INCOMPLETE fStatus = iota - 1
	FILTER_ALL
	FILTER_COMPLETE
)

type Filter struct {
	Worker int
	Addr   int
	Recent bool
	Status fStatus
}

func (e Endpoints) print(c echo.Context) error {
	d := struct {
		Data   Data
		Helper Helper
	}{
		Data:   data,
		Helper: helper,
	}

	return c.Render(200, "printtable", d)
}

func (e Endpoints) alltasks_GET(c echo.Context) error {
	fids := []int{}

	for _, v := range data.Tasks {
		fids = append(fids, v.Id)
	}

	data.Filtered = fids

	filter := Filter{}

	d := struct {
		Filter Filter
		Data   Data
		Helper Helper
	}{
		Filter: filter,
		Data:   data,
		Helper: helper,
	}

	return c.Render(200, "alltasks", d)
}

func (e Endpoints) alltasks_filter(c echo.Context) error {
	wid, _ := strconv.Atoi(c.FormValue("f-worker"))
	addr, _ := strconv.Atoi(c.FormValue("f-addr"))
	recent := c.FormValue("f-recent") == "1"
	status, _ := strconv.Atoi(c.FormValue("f-status"))

	filter := Filter{
		Worker: wid,
		Addr:   addr,
		Recent: recent,
		Status: fStatus(status),
	}

	tasks := []db.Task{}
	fids := []int{}

	cmp := func(data, form int) bool {
		return form == 0 || data == form
	}

	cut := time.Now().AddDate(0, 0, -15)

	for _, v := range data.Tasks {
		if !cmp(v.Addr_obj, addr) {
			continue
		}

		if !cmp(v.Worker, wid) {
			continue
		}

		if recent && v.Created_date.Before(cut) {
			continue
		}

		switch fStatus(status) {
		case FILTER_COMPLETE:
			if !v.Mark_date.Valid {
				continue
			}
		case FILTER_INCOMPLETE:
			if v.Mark_date.Valid {
				continue
			}
		case FILTER_ALL:
		default:
			return echo.ErrBadRequest
		}

		tasks = append(tasks, *v)
		fids = append(fids, v.Id)
	}

	data.Filtered = fids

	temp := struct {
		Tasks []db.Task
		Cats  CatMap
		Addrs map[int]*db.Address
		UWs   map[int]*UserWorker
	}{
		Tasks: tasks,
		Cats:  data.Cats,
		Addrs: data.Addrs,
		UWs:   data.UWs,
	}

	d := struct {
		Filter Filter
		Data   any
		Helper Helper
	}{
		Filter: filter,
		Data:   temp,
		Helper: helper,
	}

	return c.Render(200, "alltasks-tbody", d)
}
