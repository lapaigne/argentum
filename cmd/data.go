package main

import (
	"argentum/db"
	"fmt"
)

type Data struct {
	Incomplete []db.Task
	Workers    []db.Worker
	Addresses  []db.Address
	Categories []db.Category
	Today      string
}

// fetch mostly static data: categories, workers, adresses, etc
func FetchRare() error {
	var err error

	data.Workers, err = db.GetWorkers()
	if err != nil {
		fmt.Printf("Error on getting workers list: %s", err)
		return err
	}

	data.Addresses, err = db.GetAddresses()
	if err != nil {
		fmt.Printf("Error on getting addresses list: %s", err)
		return err
	}

	data.Categories, err = db.GetCategories()
	if err != nil {
		fmt.Printf("Error on getting categories list: %s", err)
		return err
	}

	return nil
}

func FetchIncomplete() error {
	var err error

	data.Incomplete, err = db.GetIncomplete()
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
