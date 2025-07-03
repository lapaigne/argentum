package db

import "log"

type Worker struct {
	id     int
	f_name string
	i_name string
	o_name string // this string might be empty
}

func getWorkers(id int) *Worker {

	rows, err := db.Query("SELECT * FROM public.workers")

	if err != nil {
		log.Fatal(err)
	}

	rows.Close()
	// process rows

	return nil
}
