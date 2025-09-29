package db

import (
	"database/sql"
	"fmt"
	"time"
)

type Task struct {
	Id             int
	Cat_1          int
	Cat_2          int
	Cat_3          int
	Desc           string
	Addr_obj       int
	Created_date   time.Time
	Until_date     sql.NullTime
	Mark_date      sql.NullTime
	Completed_date sql.NullTime
	Comment        string
	Worker         int
}

func AddTask(t Task) error {

	query := `
	INSERT INTO public.tasks 
	(cat_1, cat_2, cat_3, "desc", addr_obj, created_date, until_date, comment, worker)
	VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9);
	`
	_, err := db.Exec(query, t.Cat_1, t.Cat_2, t.Cat_3, t.Desc, t.Addr_obj, t.Created_date, t.Until_date, t.Comment, t.Worker)
	return err
}

// m is mark_date

// c is completed_date
func UpdateTask(id int, m sql.NullTime, c sql.NullTime) (bool, error) {

	query := `SELECT "mark_date", "completed_date" FROM public.tasks WHERE id = $1`
	var md, cd sql.NullTime

	if err := db.QueryRow(query, id).Scan(&md, &cd); err != nil {
		return false, err
	}

	if cd.Valid && c.Valid || md.Valid && m.Valid {
		return false, fmt.Errorf(`non-null task status for task %d`, id)
	}

	query = `UPDATE public.tasks SET "mark_date" = $1, "completed_date" = $2 WHERE id = $3`

	_, err := db.Exec(query, m, c, id)
	if err != nil {
		return false, err
	}

	return true, nil
}

func GetAllTasks() ([]Task, error) {

	query := "SELECT * FROM public.tasks;"
	rows, err := db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	res := []Task{}
	for rows.Next() {
		var t Task
		if err := rows.Scan(
			&t.Id,
			&t.Cat_1,
			&t.Cat_2,
			&t.Cat_3,
			&t.Desc,
			&t.Addr_obj,
			&t.Created_date,
			&t.Until_date,
			&t.Mark_date,
			&t.Completed_date,
			&t.Comment,
			&t.Worker,
		); err != nil {
			return nil, err
		}

		res = append(res, t)
	}

	if err = rows.Err(); err != nil {
		return res, err
	}

	return res, nil
}

func TasksByWorker(w int) ([]Task, error) {

	res := []Task{}
	query := "SELECT (id) FROM public.tasks WHERE 'actor' = $1"
	rows, err := db.Query(query, w)
	defer rows.Close()

	return res, err
}
