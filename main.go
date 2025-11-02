package main

import (
	"database/sql"

	_ "modernc.org/sqlite"
)

func main() {
	db, err := sql.Open("sqlite3", "./data/tasks.db")
	if err != nil {
		panic(err)
	}
	defer db.Close()

	// Your code here
}
