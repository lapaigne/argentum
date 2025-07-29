package main

import (
	"argentum/db"
	"database/sql"
	"fmt"
	"time"
)

const (
	// DD.MM.YYYY
	DateFormat = "02.01.2006"
)

type Data struct {
	Raw    RawData
	Mapped MappedData
	Helper Helper
}

type Helper struct {
	Today   string
	Cat_1   sql.NullInt32
	Cat_2   sql.NullInt32
	Addr    int
	Format  func(time.Time) string
	NFormat func(sql.NullTime) string
	Sort    TaskDisplay
}

type RawData struct {
	Tasks      []db.Task
	Workers    []db.Worker
	Addresses  []db.Address
	Categories []db.Category
}

type MappedData struct {
	ProperWorkers map[int]*db.Worker
	Workers       map[int]string
	Addresses     map[int]string
	Categories    map[int]string
	Tasks         map[int]*db.Task
	WorkersDifs   map[int]*WorkersDif
}

type WorkersDif struct {
	Before db.Worker
	After  db.Worker
	Update bool
}

func Format(t time.Time) string {

	d := t.Format(DateFormat)
	return d
}

func NFormat(t sql.NullTime) string {

	if !t.Valid {
		return "-"
	}
	d := t.Time.Format(DateFormat)
	return d
}

func (d *Data) ResetHelper() {

	d.Helper.Cat_1 = sql.NullInt32{Int32: -1, Valid: true}
	d.Helper.Cat_2 = sql.NullInt32{Int32: -1, Valid: true}

	d.Helper.Addr = -1
}

func (d *Data) Init() {

	d.Mapped.Workers = make(map[int]string)
	d.Mapped.ProperWorkers = make(map[int]*db.Worker)
	d.Mapped.Addresses = make(map[int]string)
	d.Mapped.Categories = make(map[int]string)
	d.Mapped.Tasks = make(map[int]*db.Task)

	d.Helper.Cat_1 = sql.NullInt32{Int32: -1, Valid: true}
	d.Helper.Cat_2 = sql.NullInt32{Int32: -1, Valid: true}

	d.Helper.Format = Format
	d.Helper.NFormat = NFormat

	// d.Helper.Sort = sorting
}

func (d *Data) Fill() error {

	for _, v := range d.Raw.Workers {
		d.Mapped.ProperWorkers[v.Id] = &v
	}

	for _, v := range d.Raw.Workers {
		d.Mapped.Workers[v.Id] = v.F_name + " " + v.I_name + " " + v.O_name
	}

	for _, v := range d.Raw.Addresses {
		d.Mapped.Addresses[v.Id] = v.Address
	}

	for _, v := range d.Raw.Categories {
		d.Mapped.Categories[v.Id] = v.Name
	}

	for _, v := range d.Raw.Tasks {
		d.Mapped.Tasks[v.Id] = &v
	}

	return nil
}

// fetch mostly static data: categories, workers, adresses, etc
func FetchRare() error {
	var err error

	data.Raw.Workers, err = db.GetWorkers()
	if err != nil {
		fmt.Printf("Error on getting workers list: %s", err)
		return err
	}

	data.Raw.Addresses, err = db.GetAddresses()
	if err != nil {
		fmt.Printf("Error on getting addresses list: %s", err)
		return err
	}

	data.Raw.Categories, err = db.GetCategories()
	if err != nil {
		fmt.Printf("Error on getting categories list: %s", err)
		return err
	}

	return nil
}

func FetchTasks() error {
	var err error

	data.Raw.Tasks, err = db.GetAllTasks()
	if err != nil {
		fmt.Printf("Error on fetching ALL incomplete tasks: %s", err)
		return err
	}

	return nil
}

func Fetch() error {
	var err error
	_, err = db.GetAllTasks()
	if err != nil {
		return err
	}

	return nil
}
