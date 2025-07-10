package db

import (
	"fmt"
)

type Address struct {
	Id      int
	Address string
}

// adds address objects to list in the db
func AddAddress(address string) error {

	query := fmt.Sprintf("INSERT INTO public.addr_objs (address) VALUES %s;", address)
	_, err := db.Exec(query)
	if err != nil {
		return err
	}

	return nil
}

func GetAddresses() ([]Address, error) {

	query := "SELECT * FROM public.addr_objs"
	rows, err := db.Query(query)
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

	if err = rows.Err(); err != nil {
		return res, err
	}

	return res, nil
}
