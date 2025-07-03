package db

import (
	"database/sql"
	"log"

	_ "github.com/lib/pq"
)

func OpenConn() {
	connStr := "user=lapaigne dbname=dbdev sslmode=verify-full"
	var err error
	db, err = sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}
}

func CloseConn() {
	db.Close()
}
