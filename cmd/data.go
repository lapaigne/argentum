package main

import (
	"argentum/db"
	"database/sql"
	"maps"
	"slices"
	"time"
)

type Data struct {
	Tasks map[int]*db.Task
	Cats  CatMap
	Addrs map[int]*db.Address
	UWs   map[int]*UserWorker

	Filtered []int

	MTask   int
	MCat    int
	MAddr   int
	MUser   int
	MWorker int
}

type Helper struct {
	Today    string
	Format   func(time.Time) string
	NFormat  func(sql.NullTime) string
	Status   func(sql.NullTime) string
	PNFormat func(sql.NullTime) string
	PFormat  func(time.Time) string
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

func (uw *UserWorker) SoftDelete() {
	uw.U.Level = ACC_NONE
}

type CatMap map[int]*db.Category
type CatArr []db.Category

func (cats CatMap) Sort() CatArr {
	vals := slices.Collect(maps.Values(cats))
	slices.SortFunc(vals, func(a, b *db.Category) int {
		if a.Id > b.Id {
			return 1
		} else {
			return -1
		}
	})

	sorted := CatArr{}

	l2 := CatArr{}
	l3 := CatArr{}

	for _, v := range vals {
		if v.Level == 2 {
			l2 = append(l2, *v)
		}
		if v.Level == 3 {
			l3 = append(l3, *v)
		}
	}

	for _, x := range vals {
		if x.Level == 1 {
			sorted = append(sorted, *x)

			for _, y := range l2 {
				if y.Parent.Int32 == int32(x.Id) {
					sorted = append(sorted, y)

					for _, z := range l3 {
						if z.Parent.Int32 == int32(y.Id) {
							sorted = append(sorted, z)
						}
					}

				}
			}

		}
	}

	return sorted
}
