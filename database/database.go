package database

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

var db *sql.DB

// ConnectDB connects to the SQLite database and returns a pointer to the database object.
func ConnectDB() (*sql.DB, error) {
	var err error
	db, err = sql.Open("sqlite3", "./mydb.db")
	if err != nil {
		log.Fatal(err)
		return nil, err
	}

	err = CreateTables()
	if err != nil {
		log.Fatal(err)
		return nil, err
	}

	return db, nil
}

// createTables creates tables for collections, tags, and data points in the database.
func CreateTables() error {
	collectionsTable := `
		CREATE TABLE IF NOT EXISTS collections (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			name TEXT NOT NULL
		);
	`

	tagsTable := `
		CREATE TABLE IF NOT EXISTS tags (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			name TEXT NOT NULL,
			collection_id INTEGER NOT NULL,
			FOREIGN KEY (collection_id) REFERENCES collections(id) ON DELETE CASCADE
		);
	`

	dataPointsTable := `
		CREATE TABLE IF NOT EXISTS data_points (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			content TEXT NOT NULL,
			tag_id INTEGER NOT NULL,
			FOREIGN KEY (tag_id) REFERENCES tags(id) ON DELETE CASCADE
		);
	`

	_, err := db.Exec(collectionsTable)
	if err != nil {
		return fmt.Errorf("error creating collections table: %v", err)
	}

	_, err = db.Exec(tagsTable)
	if err != nil {
		return fmt.Errorf("error creating tags table: %v", err)
	}

	_, err = db.Exec(dataPointsTable)
	if err != nil {
		return fmt.Errorf("error creating data points table: %v", err)
	}

	return nil
}

// ExecuteQuery executes the given SQL query and returns the result.
func ExecuteQuery(query string, args ...interface{}) (*sql.Rows, error) {
	rows, err := db.Query(query, args...)
	if err != nil {
		return nil, fmt.Errorf("error executing query: %v", err)
	}

	return rows, nil
}

// ExecuteNonQuery executes the given SQL query that does not return any rows.
func ExecuteNonQuery(query string, args ...interface{}) error {
	_, err := db.Exec(query, args...)
	if err != nil {
		return fmt.Errorf("error executing non-query: %v", err)
	}

	return nil
}
