package db

import (
	"context"
	"database/sql"
	"fmt"
	"time"
)

type Task struct {
	Id             int
	Cat1           int
	Cat2           int
	Cat3           int
	Desc           string
	Addr_obj       int
	Created_date   time.Time
	Until_date     sql.NullTime
	Mark_date      sql.NullTime
	Completed_date sql.NullTime
	Comment        string
	Worker         int
}

func AddTask(t Task, ctx context.Context) (int, error) {
	stmt, err := db.Prepare(`INSERT INTO public.tasks ("cat_1", "cat_2", "cat_3", "desc", "addr_obj", "created_date", "until_date", "comment", "worker") VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9); `)
	if err != nil {
		return -1, err
	}
	defer stmt.Close()

	res, err := stmt.ExecContext(ctx, t.Cat1, t.Cat2, t.Cat3, t.Desc, t.Addr_obj, t.Created_date, t.Until_date, t.Comment, t.Worker)
	if err != nil {
		return -1, err
	}

	id, err := res.LastInsertId()
	return int(id), err
}

// m is mark_date
//
// c is completed_date
func UpdateTask(id int, m sql.NullTime, c sql.NullTime) error {
	selStmt, err := db.Prepare(`SELECT "mark_date", "completed_date" FROM public.tasks WHERE id = $1`)
	if err != nil {
		return err
	}
	defer selStmt.Close()

	var md, cd sql.NullTime

	if err := selStmt.QueryRow(id).Scan(&md, &cd); err != nil {
		return err
	}

	if cd.Valid && c.Valid || md.Valid && m.Valid {
		return fmt.Errorf(`non-null task status for task %d`, id)
	}

	updStmt, err := db.Prepare(`UPDATE public.tasks SET "mark_date" = $1, "completed_date" = $2 WHERE id = $3`)
	if err != nil {
		return err
	}
	defer updStmt.Close()

	_, err = updStmt.Exec(m, c, id)
	return err
}

func GetTasks(ctx context.Context) ([]Task, error) {
	res := []Task{}
	stmt, err := db.PrepareContext(ctx, `SELECT * FROM public.tasks;`)
	if err != nil {
		return res, err
	}
	defer stmt.Close()

	rows, err := stmt.QueryContext(ctx)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var t Task
		if err := rows.Scan(
			&t.Id,
			&t.Cat1,
			&t.Cat2,
			&t.Cat3,
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

// unused
func TasksByWorker(w int) ([]Task, error) {
	res := []Task{}
	query := "SELECT (id) FROM public.tasks WHERE 'actor' = $1"
	rows, err := db.Query(query, w)
	defer rows.Close()

	return res, err
}
