package main

import (
	"argentum/db"
	"database/sql"
	"fmt"
)

type Data struct {
	Raw    RawData
	Mapped MappedData
	Helper HelperData
}

type HelperData struct {
	Today string
	Cat_1 sql.NullInt32
	Cat_2 sql.NullInt32
}

type RawData struct {
	Tasks      []db.Task
	Workers    []db.Worker
	Addresses  []db.Address
	Categories []db.Category
}

type MappedData struct {
	Workers    map[int]string
	Addresses  map[int]string
	Categories map[int]string
}

func (d *Data) Init() {

	d.Mapped.Workers = make(map[int]string)
	d.Mapped.Addresses = make(map[int]string)
	d.Mapped.Categories = make(map[int]string)
}

func (d *Data) Fill() error {

	for _, v := range d.Raw.Workers {
		d.Mapped.Workers[v.Id] = v.F_name + " " + v.I_name + " " + v.O_name
	}

	for _, v := range d.Raw.Addresses {
		d.Mapped.Addresses[v.Id] = v.Address
	}

	for _, v := range d.Raw.Categories {
		d.Mapped.Categories[v.Id] = v.Name
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

	data.Raw.Tasks, err = db.GetTasks()
	if err != nil {
		fmt.Printf("Error on fetching ALL incomplete tasks: %s", err)
		return err
	}

	return nil
}

func Fetch() error {
	var err error
	_, err = db.GetTasks()
	if err != nil {
		return err
	}
	return nil

}
