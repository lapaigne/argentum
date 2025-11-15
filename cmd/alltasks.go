package main

import (
	"argentum/db"
	"fmt"
	"net/http"
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
	filter := Filter{}

	if clw, _ := c.Cookie("flw"); clw != nil {
		filter.Worker, _ = strconv.Atoi(clw.Value)
	}

	if cla, _ := c.Cookie("fla"); cla != nil {
		filter.Addr, _ = strconv.Atoi(cla.Value)
	}

	if cre, _ := c.Cookie("fre"); cre != nil {
		filter.Recent, _ = strconv.ParseBool(cre.Value)
	}

	if cst, _ := c.Cookie("fst"); cst != nil {
		s, _ := strconv.Atoi(cst.Value)
		filter.Status = fStatus(s)
	}

	tasks, err := at_filter(filter)
	if err != nil {
		fmt.Println(len(tasks))
	}

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
		Data   any
		Helper Helper
	}{
		Data:   temp,
		Helper: helper,
	}

	return c.Render(200, "printtable", d)
}

func (e Endpoints) alltasks_GET(c echo.Context) error {
	filter := Filter{}

	if clw, _ := c.Cookie("flw"); clw != nil {
		filter.Worker, _ = strconv.Atoi(clw.Value)
	}

	if cla, _ := c.Cookie("fla"); cla != nil {
		filter.Addr, _ = strconv.Atoi(cla.Value)
	}

	tasks, err := at_filter(filter)
	if err != nil {
		fmt.Println(len(tasks))
	}

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

func (e Endpoints) alltasks_filter(c echo.Context) error {
	wid, _ := strconv.Atoi(c.FormValue("f-worker"))
	addr, _ := strconv.Atoi(c.FormValue("f-addr"))
	recent := c.FormValue("f-recent") == "1"
	status, _ := strconv.Atoi(c.FormValue("f-status"))

	// worker
	c.SetCookie(&http.Cookie{
		Name:  "flw",
		Value: strconv.Itoa(wid),
		Path:  "/",
	})

	// addr
	c.SetCookie(&http.Cookie{
		Name:  "fla",
		Value: strconv.Itoa(addr),
		Path:  "/",
	})

	// recent
	c.SetCookie(&http.Cookie{
		Name:  "fre",
		Value: strconv.FormatBool(recent),
		Path:  "/",
	})

	// status
	c.SetCookie(&http.Cookie{
		Name:  "fst",
		Value: strconv.Itoa(int(status)),
		Path:  "/",
	})

	filter := Filter{
		Worker: wid,
		Addr:   addr,
		Recent: recent,
		Status: fStatus(status),
	}

	tasks, err := at_filter(filter)
	if err != nil {
		return err
	}

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

func at_filter(filter Filter) ([]db.Task, error) {
	worker, addr, recent, status := filter.Worker, filter.Addr, filter.Recent, filter.Status
	tasks := []db.Task{}

	cmp := func(d, form int) bool {
		return form == 0 || d == form
	}

	cut := time.Now().AddDate(0, 0, -15)

	for _, v := range data.Tasks {
		if !cmp(v.Addr_obj, addr) {
			continue
		}

		if !cmp(v.Worker, worker) {
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
			return tasks, fmt.Errorf("invalid filter value")
		}

		tasks = append(tasks, *v)
	}

	return tasks, nil
}
