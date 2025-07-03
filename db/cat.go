package db

import (
	"log"
)

type Cat struct {
	current int
	parent  int
}

func getCats() []Cat {
	res := []Cat{}

	// NOTE: decide to either read data from db (if it is stored there)
	// or read data from local file
	// and decide on the file format

	_, err := db.Query("*")
	if err != nil {
		log.Fatal(err)
	}

	return res
}
