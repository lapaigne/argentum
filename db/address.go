package db

import (
	"context"
)

type Address struct {
	Id      int
	Address string
}

// adds address objects to list in the db
func AddAddress(address string) (int, error) {
	stmt, err := db.Prepare(`INSERT INTO public.addr_objs ("address") VALUES $1 RETURNING id`)
	if err != nil {
		return -1, err
	}
	defer stmt.Close()

	var id *int
	if err := stmt.QueryRow(address).Scan(&id); err != nil {
		return -1, err
	}

	return int(*id), nil
}

func GetAddresses(ctx context.Context) ([]Address, error) {
	stmt, err := db.PrepareContext(ctx, `SELECT * FROM public.addr_objs ORDER BY id ASC`)
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
		if err := rows.Scan(&address.Id, &address.Address); err != nil {
			return res, err
		}
		res = append(res, address)
	}

	return res, rows.Err()
}
