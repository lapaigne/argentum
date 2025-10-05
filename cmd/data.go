package main

import (
	"argentum/db"
	"database/sql"
	"time"
)

type Data struct {
	Tasks map[int]db.Task
	Cats  map[int]db.Category
	Addrs map[int]db.Address
	UWs   map[int]UserWorker

	MTask int
	MCat  int
	MAddr int
}

type Helper struct {
	Today   string
	Format  func(time.Time) string
	NFormat func(sql.NullTime) string
	Status  func(sql.NullTime) string
	Levels  *map[string]int
}

type UserWorker struct {
	Id  int
	W   db.Worker
	U   db.User
	DW  bool // if true, update worker
	DU  bool // if true, update user
	New bool // if true, add new user and new worker
}
