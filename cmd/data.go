package main

import (
	"argentum/db"
	"database/sql"
	"time"
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

type Helper struct {
	Today   string
	Format  func(time.Time) string
	NFormat func(sql.NullTime) string
	Levels  *map[string]int
}

type AForm struct {
	Cat_1 sql.NullInt32
	Cat_2 sql.NullInt32
	Addr  int
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

func Format(t time.Time) string {
	return t.Format(time.DateOnly)
}

func NFormat(t sql.NullTime) string {
	if !t.Valid {
		return ""
	}

	return t.Time.Format(time.DateOnly)
}
