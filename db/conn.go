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
	// TODO: use proper <USER>, <PASSWORD>, <HOST>
	connStr := fmt.Sprintf("user=lapaigne password=%s dbname=aqua host=laptopaigne sslmode=verify-full", os.Getenv("PSQLPASS"))
	var err error
	db, err = sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}
}

func CloseConn() {
	db.Close()
}
