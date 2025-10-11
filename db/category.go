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
	Active bool
}

func UpdateCategory(id int, cat Category) error {
	stmt, err := db.Prepare(`UPDATE public.categories SET "name" = $1, "active" = $2 WHERE "id" = $3`)
	if err != nil {
		return err
	}
	defer stmt.Close()
	_, err = stmt.Exec(cat.Name, cat.Active, id)
	return err
}

func AddCategory(cat Category) (int, error) {
	if cat.Level > 3 {
		s := fmt.Sprintf("Invalid category level of %d, when adding category <%s>", cat.Level, cat.Name)
		return -1, errors.New(s)
	}

	stmt, err := db.Prepare(`INSERT INTO public.categories ("parent_id", "name", "level", "active") VALUES ($1, $2, $3, $4) RETURNING id`)
	if err != nil {
		return -1, err
	}
	defer stmt.Close()

	var id *int
	if err := stmt.QueryRow(cat.Parent, cat.Name, cat.Level, cat.Active).Scan(&id); err != nil {
		return -1, err
	}

	return int(*id), err
}

func GetCategories(ctx context.Context) ([]Category, error) {
	stmt, err := db.PrepareContext(ctx, `SELECT "id", "parent_id", "name", "level", "active" FROM public.categories ORDER BY id ASC`)
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
		if err := rows.Scan(&cat.Id, &cat.Parent, &cat.Name, &cat.Level, &cat.Active); err != nil {
			return res, err
		}
		res = append(res, cat)
	}

	return res, rows.Err()
}
