package main

import (
	"argentum/db"
	"strconv"

	"github.com/labstack/echo/v4"
)

type Filter struct {
	Worker int
	Addr   int
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

	filter := Filter{
		Worker: 0,
		Addr:   0,
	}

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

	filter := Filter{
		Worker: wid,
		Addr:   addr,
	}

	tasks := []db.Task{}
	fids := []int{}

	cmp := func(data, form int) bool {
		if form == 0 {
			return true
		}
		if data != form {
			return false
		}
		return true
	}

	for _, v := range data.Tasks {
		if cmp(v.Worker, wid) && cmp(v.Addr_obj, addr) {
			tasks = append(tasks, *v)
			fids = append(fids, v.Id)
		}
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

	return c.Render(200, "alltasks", d)
}
