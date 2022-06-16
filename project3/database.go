package main

import (
	"database/sql"
	"fmt"

	log "github.com/sirupsen/logrus"

	_ "github.com/mattn/go-sqlite3"
)

var db *sql.DB

func initDB() {
	var err error
	db, err = sql.Open("sqlite3", "test.db")
	if err != nil {
		log.Fatal(err)
	}
}

// Create the tables in our database
func createTables() {

	////////////////////////////
	// BEGIN: YOUR CODE HERE
	////////////////////////////
	statement := `
		CREATE TABLE IF NOT EXISTS users (id INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
							username TEXT,
							password TEXT,
							salt TEXT
							);
		CREATE TABLE IF NOT EXISTS sessions (id INTEGER NOT NULL PRIMARY KEY,
							   username TEXT,
							   token TEXT,
							   expires INTEGER
							   );
		CREATE TABLE IF NOT EXISTS files (id INTEGER NOT NULL PRIMARY KEY,
							yourfield TEXT
							);`
	// TODO: modify the schema of the files table to help implement tasks 3-6.
	// do NOT modify the schema of the sessions or users tables.

	// NOTE: You need to delete the test.db file after making changes to the SQL above.
	// Otherwise your schema will not be applied.

	////////////////////////////
	// END: YOUR CODE HERE
	////////////////////////////

	log.Info("setting up database")
	_, err := db.Exec(statement)
	if err != nil {
		log.Fatal(err)
	}
	log.Info("setting up database done")
}

// Remove all tables from the database
func dropTables() {
	log.Printf("dropping all tables")
	tables := []string{"users", "sessions", "files"}
	for _, table := range tables {
		_, err := db.Exec("DROP TABLE " + table)
		if err != nil {
			log.Fatal(err)
		}
	}
}

// Print the contents of a table (helpful for debugging)
// Code from https://github.com/Go-SQL-Driver/MySQL/wiki/Examples#rawbytes
func printTable(db *sql.DB, table string) {
	query := fmt.Sprintf("SELECT * FROM %s", table)
	rows, err := db.Query(query)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	// Get column names
	columns, err := rows.Columns()
	if err != nil {
		panic(err.Error()) // proper error handling instead of panic in your app
	}

	// Make a slice for the values
	values := make([]sql.RawBytes, len(columns))

	// rows.Scan wants '[]interface{}' as an argument, so we must copy the
	// references into such a slice
	// See http://code.google.com/p/go-wiki/wiki/InterfaceSlice for details
	scanArgs := make([]interface{}, len(values))
	for i := range values {
		scanArgs[i] = &values[i]
	}

	// Fetch rows
	for rows.Next() {
		// get RawBytes from data
		err = rows.Scan(scanArgs...)
		if err != nil {
			panic(err.Error())
		}

		// Here we just print each column as a string.
		var value string
		for i, col := range values {
			// Here we can check if the value is nil (NULL value)
			if col == nil {
				value = "NULL"
			} else {
				value = string(col)
			}
			fmt.Println(columns[i], ": ", value)
		}
		fmt.Println("-----------------------------------")
	}
	if err = rows.Err(); err != nil {
		panic(err.Error()) // proper error handling instead of panic in your app
	}
}
