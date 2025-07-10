package db

import (
	"fmt"
)

type Worker struct {
	Id     int
	F_name string
	I_name string
	O_name string // this string might be empty
}

func AddWorker(f, i, o string) error {

	query := fmt.Sprintf("INSERT INTO public.workers (f_name, i_name, o_name) VALUES ('%s', '%s', '%s');", f, i, o)
	_, err := db.Exec(query)
	if err != nil {
		return err
	}

	return nil
}

func GetWorker(id int) (Worker, error) {

	var worker Worker
	err := db.QueryRow("SELECT id, f_name, i_name, o_name FROM public.workers WHERE id = $1", id).
		Scan(&worker.Id, &worker.F_name, &worker.I_name, &worker.O_name)
	if err != nil {
		return Worker{}, err
	}

	return worker, nil
}

func GetWorkers() ([]Worker, error) {

	rows, err := db.Query("SELECT id, f_name, i_name, o_name FROM public.workers")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var res []Worker
	for rows.Next() {
		var w Worker
		if err := rows.Scan(&w.Id, &w.F_name, &w.I_name, &w.O_name); err != nil {
			return res, err
		}
		res = append(res, w)
	}

	if err = rows.Err(); err != nil {
		return res, err
	}

	return res, nil
}
