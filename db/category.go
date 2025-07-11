package db

import (
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

func AddCategory(parent, level int, name string) error {

	if level > 3 {
		s := fmt.Sprintf("Invalid category level of %d, when adding category <%s>", level, name)
		return errors.New(s)
	}

	query := fmt.Sprintf("INSERT INTO public.categories (parent_id, name, level) VALUES (%d, %s, %d);", parent, name, level)
	_, err := db.Exec(query)
	if err != nil {
		return err
	}

	return nil

}

func GetCategories() ([]Category, error) {

	query := "SELECT * FROM public.categories"
	rows, err := db.Query(query)
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

	if err = rows.Err(); err != nil {
		return res, err
	}

	return res, nil
}
