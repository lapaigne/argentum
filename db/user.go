package db

import (
	"context"
	"database/sql"
)

type User struct {
	Id     int
	Worker int
	Login  string
	Hash   string
	Level  int
	Token  sql.NullString
}

var AccessLevels = map[string]int{
	"admin":      100,
	"dispatcher": 50,
	"worker":     10,
}

func UpdateToken(refresh string, id int) error {
	stmt, err := db.Prepare(`UPDATE public.users SET "refresh" = $1 WHERE "id" = $2`)
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(refresh, id)
	return err
}

func ValidateToken(refresh string, id int) error {
	stmt, err := db.Prepare(`SELECT 1 FROM public.users WHERE "refresh" = $1 AND "id" = $2`)
	if err != nil {
		return err
	}
	defer stmt.Close()

	var a int
	return stmt.QueryRow(refresh, id).Scan(a)
}

func GetUser(login string) (User, error) {
	stmt, err := db.Prepare(`SELECT * FROM public.users WHERE "login" = $1`)
	if err != nil {
		return User{}, err
	}
	defer stmt.Close()

	var u User
	err = stmt.QueryRow(login).Scan(&u.Id, &u.Worker, &u.Login, &u.Hash, &u.Level, &u.Token)
	return u, nil
}

func AddUser(u User) (int, error) {
	stmt, err := db.Prepare(`INSERT INTO public.users ("worker", "login", "hash", "level") VALUES ($1, $2, $3, $4)`)
	if err != nil {
		return -1, err
	}
	defer stmt.Close()

	res, err := stmt.Exec(u.Worker, u.Login, u.Hash, u.Level)
	if err != nil {
		return -1, err
	}

	id, err := res.LastInsertId()
	return int(id), err
}

func GetUsers(ctx context.Context) ([]User, error) {
	stmt, err := db.PrepareContext(ctx, `SELECT "id", "worker", "login", "level" FROM public.users`)
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
