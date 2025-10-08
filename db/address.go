package db

import (
	"context"
)

type Address struct {
	Id      int
	Address string
	Active  bool
}

func UpdateAddress(id int, a Address) error {
	stmt, err := db.Prepare(`UPDATE public.addr_objs SET "address" = $1, "active" = $2 WHERE "id" = $3`)
	if err != nil {
		return err
	}
	defer stmt.Close()
	_, err = stmt.Exec(a.Address, a.Active, id)
	return err
}

func AddAddress(a Address) (int, error) {
	stmt, err := db.Prepare(`INSERT INTO public.addr_objs ("address", "active") VALUES ($1, $2) RETURNING id`)
	if err != nil {
		return -1, err
	}
	defer stmt.Close()

	var id *int
	if err := stmt.QueryRow(a.Address, a.Active).Scan(&id); err != nil {
		return -1, err
	}

	return int(*id), nil
}

func GetAddresses(ctx context.Context) ([]Address, error) {
	stmt, err := db.PrepareContext(ctx, `SELECT "id", "address", "active" FROM public.addr_objs ORDER BY id ASC`)
	if err != nil {
		return []Address{}, err
	}
	defer stmt.Close()

	rows, err := stmt.QueryContext(ctx)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	res := []Address{}
	for rows.Next() {
		var address Address
		if err := rows.Scan(&address.Id, &address.Address, &address.Active); err != nil {
			return res, err
		}
		res = append(res, address)
	}

	return res, rows.Err()
}
