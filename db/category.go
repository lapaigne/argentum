package db

import (
	"fmt"
)

type Category struct {
	Id     int
	Parent int
	Name   string
}

func AddCategory(parent int, name string) error {

	query := fmt.Sprintf("INSERT INTO public.categories (parent_id, name) VALUES (%d, %s);", parent, name)
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
		var parent *int
		if err := rows.Scan(&cat.Id, &parent, &cat.Name); err != nil {
			return res, err
		}
		if parent == nil {
			cat.Parent = -1
		} else {
			cat.Parent = *parent
		}
		res = append(res, cat)
	}

	if err = rows.Err(); err != nil {
		return res, err
	}

	return res, nil
}
