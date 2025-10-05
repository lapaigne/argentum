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
	stmt, err := db.Prepare(`INSERT INTO public.addr_objs ("address") VALUES $1;`)
	if err != nil {
		return -1, err
	}
	defer stmt.Close()

	res, err := stmt.Exec(address)
	if err != nil {
		return -1, err
	}

	id, err := res.LastInsertId()
	return int(id), nil
}

func GetAddresses(ctx context.Context) ([]Address, error) {
	stmt, err := db.PrepareContext(ctx, `SELECT * FROM public.addr_objs`)
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
