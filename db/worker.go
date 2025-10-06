package db

import (
	"context"
)

type Worker struct {
	Id     int
	F_name string
	I_name string
	O_name string
}

func UpdateWorker(id int, w Worker) error {
	stmt, err := db.Prepare(`UPDATE public.workers SET "f_name" = $1, "i_name" = $2, "o_name" = $3 WHERE "id" = $4`)
	if err != nil {
		return err
	}
	defer stmt.Close()
	_, err = stmt.Exec(w.F_name, w.I_name, w.O_name, id)
	return err
}

func AddWorker(f, i, o string) (int, error) {
	stmt, err := db.Prepare(`INSERT INTO public.workers ("f_name", "i_name", "o_name") VALUES ($1, $2, $3) RETURNING id`)
	if err != nil {
		return -1, err
	}
	defer stmt.Close()

	var id *int
	if err := stmt.QueryRow(f, i, o).Scan(&id); err != nil {
		return -1, err
	}

	return int(*id), err
}

// unused
func GetWorker(id int) (Worker, error) {
	stmt, err := db.Prepare(`SELECT "id", "f_name", "i_name", "o_name" FROM public.workers WHERE "id" = $1`)
	if err != nil {
		return Worker{}, err
	}
	defer stmt.Close()

	var worker Worker
	if err = stmt.QueryRow(id).Scan(&worker.Id, &worker.F_name, &worker.I_name, &worker.O_name); err != nil {
		return Worker{}, err
	}

	return worker, nil
}

func GetWorkers(ctx context.Context) ([]Worker, error) {
	stmt, err := db.PrepareContext(ctx, `SELECT "id", "f_name", "i_name", "o_name" FROM public.workers`)
	if err != nil {
		return []Worker{}, err
	}
	defer stmt.Close()

	rows, err := stmt.QueryContext(ctx)
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
