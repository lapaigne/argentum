package main

import (
	"argentum/db"
)

type Data struct {
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
	Id  int
	W   db.Worker
	U   db.User
	DW  bool // if true, update worker
	DU  bool // if true, update user
	New bool // if true, add new user and new worker
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

	d.Workers = nil
	d.Workers = make(map[int]db.Worker)

	for _, v := range workers {
		d.Workers[v.Id] = v
	}
	if len(workers) == 0 {
		d.MWorker = 0
	} else {
		d.MWorker = workers[len(workers)-1].Id
	}

	d.UWs = nil
	d.UWs = make(map[int]UserWorker)

	for _, v := range workers {
		d.UWs[v.Id] = UserWorker{Id: v.Id, W: v, U: um[v.Id]}
	}

	tasks, err := db.GetAllTasks()
	if err != nil {
		return err
	}
	d.Tasks = nil
	d.Tasks = make(map[int]db.Task)
	for _, v := range tasks {
		d.Tasks[v.Id] = v
	}
	if len(tasks) == 0 {
		d.MTask = 0
	} else {
		d.MTask = tasks[len(tasks)-1].Id
	}

	cats, err := db.GetCategories()
	if err != nil {
		return err
	}
	d.Cats = nil
	d.Cats = make(map[int]db.Category)
	for _, v := range cats {
		d.Cats[v.Id] = v
	}
	if len(cats) == 0 {
		d.MCat = 0
	} else {
		d.MCat = cats[len(cats)-1].Id
	}

	addrs, err := db.GetAddresses()
	if err != nil {
		return err
	}
	d.Addrs = nil
	d.Addrs = make(map[int]db.Address)
	for _, v := range addrs {
		d.Addrs[v.Id] = v
	}
	if len(addrs) == 0 {
		d.MAddr = 0
	} else {
		d.MAddr = addrs[len(addrs)-1].Id
	}

	return nil
}
