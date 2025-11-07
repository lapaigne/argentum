package db

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/lib/pq"
)

var db *sql.DB

func OpenConn() {
	connStr := fmt.Sprintf("user=%s password=%s dbname=%s host=%s sslmode=require",
		os.Getenv("DB_USER"),
		os.Getenv("DB_PSWD"),
		os.Getenv("DB_NAME"),
		os.Getenv("DB_HOST"),
	)

	var err error
	db, err = sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}
}

func CloseConn() {
	db.Close()
}
