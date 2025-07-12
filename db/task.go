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
	Completed_date sql.NullTime
	Comment        string
}

type TaskWorker struct {
	Id        int
	Task_id   int
	Worker_id int
}

func AddTask(t Task) error {

	query := `
	INSERT INTO public.tasks 
	(cat_1, cat_2, cat_3, "desc", addr_obj, created_date, until_date, comment)
	VALUES ($1, $2, $3, $4, $5, $6, $7, $8);
	`
	_, err := db.Exec(query, t.Cat_1, t.Cat_2, t.Cat_3, t.Desc, t.Addr_obj, t.Created_date, t.Until_date, t.Comment)
	return err

}

func GetTasks() ([]Task, error) {

	query := "SELECT * FROM public.tasks;"
	rows, err := db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("Failed query")
	}
	defer rows.Close()

	res := []Task{}
	for rows.Next() {
		var t Task
		if err := rows.Scan(&t.Id, &t.Cat_1, &t.Cat_2, &t.Cat_3, &t.Desc, &t.Addr_obj, &t.Created_date, &t.Until_date, &t.Completed_date, &t.Comment); err != nil {
			return nil, err

		}
		res = append(res, t)
	}

	if err = rows.Err(); err != nil {
		return res, err
	}

	return res, nil
}

func TasksByWorker(id int) ([]Task, error) {

	res := []Task{}
	query := "SELECT "
	rows, err := db.Query(query)
	rows.Close()
	if err != nil {
		return res, fmt.Errorf("Failed query")
	}

	return res, nil
}
