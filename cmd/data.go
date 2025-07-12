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

func (d *MappedData) Init() {

	d.Workers = make(map[int]string)
	d.Addresses = make(map[int]string)
	d.Categories = make(map[int]string)
}

func (d *MappedData) Fill(r *RawData) error {

	for _, v := range r.Workers {
		d.Workers[v.Id] = v.F_name + " " + v.I_name + " " + v.O_name
	}

	for _, v := range r.Addresses {
		d.Workers[v.Id] = v.Address
	}

	for _, v := range r.Categories {
		d.Workers[v.Id] = v.Name
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

func FetchIncomplete() error {
	var err error

	data.Raw.Tasks, err = db.GetIncomplete()
	if err != nil {
		fmt.Printf("Error on fetching ALL incomplete tasks: %s", err)
		return err
	}

	return nil
}

func Fetch() error {
	var err error
	_, err = db.GetIncomplete()
	if err != nil {
		return err
	}
	return nil

}
