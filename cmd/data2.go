package main

import (
	"argentum/db"
)

const (
	Added = iota
	Deleted
	Edited
)

type Data struct {
	Server Mapped
	Local  Mapped
	Diffs  DiffData
}

type DiffData struct {
	Workers map[int]int
	Tasks   map[int]int
	Cats    map[int]int
	Addrs   map[int]int
}

type Mapped struct {
	Workers map[int]db.Worker
	Tasks   map[int]db.Task
	Cats    map[int]db.Category
	Addrs   map[int]db.Address
	UWs     map[int]UserWorker

	MWorker int
	MTask   int
	MCat    int
	MAddr   int
	MUser   int
}

type UserWorker struct {
	Id int
	W  db.Worker
	U  db.User
}

type RawData struct {
	Workers []db.Worker
	Tasks   []db.Task
	Cats    []db.Category
	Addrs   []db.Address
	Users   []db.User
}

func (d *Data) Fetch() error {

	workers, err := db.GetWorkers()
	if err != nil {
		return err
	}

	users, err := db.GetUsers()
	if err != nil {
		return err
	}

	um := make(map[int]db.User)
	for _, v := range users {
		um[v.Worker] = v
	}

	d.Server.Workers = nil
	d.Server.Workers = make(map[int]db.Worker)

	for _, v := range workers {
		d.Server.Workers[v.Id] = v
	}
	if len(workers) == 0 {
		d.Server.MWorker = 0
	} else {
		d.Server.MWorker = workers[len(workers)-1].Id
	}

	d.Server.UWs = nil
	d.Server.UWs = make(map[int]UserWorker)

	for _, v := range workers {
		d.Server.UWs[v.Id] = UserWorker{Id: v.Id, W: v, U: um[v.Id]}
	}

	tasks, err := db.GetAllTasks()
	if err != nil {
		return err
	}
	d.Server.Tasks = nil
	d.Server.Tasks = make(map[int]db.Task)
	for _, v := range tasks {
		d.Server.Tasks[v.Id] = v
	}
	if len(tasks) == 0 {
		d.Server.MTask = 0
	} else {
		d.Server.MTask = tasks[len(tasks)-1].Id
	}

	cats, err := db.GetCategories()
	if err != nil {
		return err
	}
	d.Server.Cats = nil
	d.Server.Cats = make(map[int]db.Category)
	for _, v := range cats {
		d.Server.Cats[v.Id] = v
	}
	if len(cats) == 0 {
		d.Server.MCat = 0
	} else {
		d.Server.MCat = cats[len(cats)-1].Id
	}

	addrs, err := db.GetAddresses()
	if err != nil {
		return err
	}
	d.Server.Addrs = nil
	d.Server.Addrs = make(map[int]db.Address)
	for _, v := range addrs {
		d.Server.Addrs[v.Id] = v
	}
	if len(addrs) == 0 {
		d.Server.MAddr = 0
	} else {
		d.Server.MAddr = addrs[len(addrs)-1].Id
	}

	return nil
}

func (d *Data) CalcDiffs() {

}

func (d *Data) WorkerDiffs() {

	for _, v := range d.Local.Workers {
		if _, ok := d.Server.Workers[v.Id]; !ok {
			d.Diffs.Workers[v.Id] = Added
		}
	}

}

func (d *Data) UpdStatus() {

}

func (d *Data) ResetLocal() {

	d.Local = d.Server
	d.CalcDiffs()

}
