package db

import (
	"context"
	"fmt"
)

type User struct {
	Id     int
	Worker int
	Login  string
	Hash   string
	Level  int
}

func DeleteUser(id int) error {
	stmt, err := db.Prepare(`UPDATE public.users SET "login" = '', "hash" = '', "level" = 0 WHERE "id" = $1`)
	if err != nil {
		return err
	}
	defer stmt.Close()
	res, err := stmt.Exec(id)
	rws, err := res.RowsAffected()
	if rws == 0 {
		return fmt.Errorf("no user with id %d", id)
	}
	return err
}

func UpdateUserLevel(id, lvl int) error {
	stmt, err := db.Prepare(`UPDATE public.users SET "level" = $1 WHERE "id" = $2`)
	if err != nil {
		return err
	}
	defer stmt.Close()
	res, err := stmt.Exec(lvl, id)
	rws, err := res.RowsAffected()
	if rws == 0 {
		return fmt.Errorf("no user with id %d", id)
	}
	return err
}

func GetUser(login string) (User, error) {
	stmt, err := db.Prepare(`SELECT "id", "worker", "hash", "level" FROM public.users WHERE "login" = $1`)
	if err != nil {
		return User{}, err
	}
	defer stmt.Close()

	var u User
	err = stmt.QueryRow(login).Scan(&u.Id, &u.Worker, &u.Hash, &u.Level)
	u.Login = login
	return u, nil
}

func AddUser(u User) (int, error) {
	stmt, err := db.Prepare(`INSERT INTO public.users ("worker", "login", "hash", "level") VALUES ($1, $2, $3, $4) RETURNING id`)
	if err != nil {
		return -1, err
	}
	defer stmt.Close()

	var id *int
	if err := stmt.QueryRow(u.Worker, u.Login, u.Hash, u.Level).Scan(&id); err != nil {
		return -1, err
	}

	return int(*id), err
}

func GetUsers(ctx context.Context) ([]User, error) {
	stmt, err := db.PrepareContext(ctx, `SELECT "id", "worker", "login", "level" FROM public.users ORDER BY id ASC`)
	if err != nil {
		return []User{}, err
	}
	defer stmt.Close()

	rows, err := stmt.QueryContext(ctx)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var res []User
	for rows.Next() {
		var u User
		if err := rows.Scan(&u.Id, &u.Worker, &u.Login, &u.Level); err != nil {
			return res, err
		}
		res = append(res, u)
	}

	if err = rows.Err(); err != nil {
		return res, err
	}

	return res, nil
}
