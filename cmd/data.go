package main

import (
	"argentum/db"
	"database/sql"
	"time"
)

type Data struct {
	Tasks map[int]*db.Task
	Cats  map[int]*db.Category
	Addrs map[int]*db.Address
	UWs   map[int]*UserWorker

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
	Id int
	W  db.Worker
	U  db.User
}

func (dest *UserWorker) SoftReplace(upd UserWorker) {
	dest.W.F_name = upd.W.F_name
	dest.W.I_name = upd.W.I_name
	dest.W.O_name = upd.W.O_name
	dest.U.Level = upd.U.Level
}
