package db

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
)

type Category struct {
	Id     int
	Parent sql.NullInt32
	Name   string
	Level  int
}

func AddCategory(parent, level int, name string) (int, error) {
	if level > 3 {
		s := fmt.Sprintf("Invalid category level of %d, when adding category <%s>", level, name)
		return -1, errors.New(s)
	}

	stmt, err := db.Prepare(`INSERT INTO public.categories ("parent_id", "name", "level") VALUES ($1, $2, $3) RETURNING id`)
	if err != nil {
		return -1, err
	}
	defer stmt.Close()

	var id *int
	if err := stmt.QueryRow(parent, name, level).Scan(&id); err != nil {
		return -1, err
	}

	return int(*id), err
}

func GetCategories(ctx context.Context) ([]Category, error) {
	stmt, err := db.PrepareContext(ctx, `SELECT * FROM public.categories`)
	if err != nil {
		return []Category{}, err
	}
	defer stmt.Close()

	rows, err := stmt.QueryContext(ctx)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	res := []Category{}
	for rows.Next() {
		var cat Category
		if err := rows.Scan(&cat.Id, &cat.Parent, &cat.Name, &cat.Level); err != nil {
			return res, err
		}
		res = append(res, cat)
	}

	return res, rows.Err()
}
