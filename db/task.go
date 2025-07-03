package db

import (
	"database/sql"
	"fmt"
	"time"
)

type Task struct {
	id         int
	cat_1      int
	cat_2      int
	cat_3      int
	desc       string
	addr_obj   int
	created_at time.Time
}

func addTask(task Task) error {
	query := ""
	_, err := db.Exec(query)
	return err
}

func incompleteTasks() ([]Task, error) {

	res := []Task{}
	query := ""
	rows, err := db.Query(query)
	rows.Close()
	if err != nil {
		return res, fmt.Errorf("Failed query")
	}

	return res, nil
}

func workerTasks(id int) ([]Task, error) {
	res := []Task{}
	query := ""
	rows, err := db.Query(query)
	rows.Close()
	if err != nil {
		return res, fmt.Errorf("Failed query")
	}

	return res, nil
}
