package db

import (
	"database/sql"
	"fmt"
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

	query := "UPDATE public.users SET refresh = $1 WHERE id = $2"
	_, err := db.Exec(query, refresh, id)
	return err
}

func ValidateToken(refresh string, id int) error {

	query := "SELECT 1 FROM public.users WHERE refresh = $1 AND id = $2"
	var a int
	err := db.QueryRow(query, refresh, id).Scan(&a)
	return err
}

func GetUser(login string) (User, error) {

	var u User
	query := "SELECT * FROM public.users WHERE login = $1"
	err := db.QueryRow(query, login).Scan(&u.Id, &u.Worker, &u.Login, &u.Hash, &u.Level, &u.Token)
	if err != nil {
		fmt.Println(err)
		return User{}, err
	}

	return u, nil
}

func AddUser(u User) error {

	query := "INSERT INTO public.users (worker, login, hash, level) VALUES ($1, $2, $3, $4)"

	_, err := db.Exec(query, u.Worker, u.Login, u.Hash, u.Level)
	if err != nil {
		return err
	}

	return nil
}
