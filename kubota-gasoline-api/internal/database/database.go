package database

import (
	"database/sql"
	"log"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

func Connect(databaseURL string) (*sql.DB, error) {
	db, err := sql.Open("sqlite3", databaseURL)
	if err != nil {
		return nil, err
	}

	db.SetMaxOpenConns(1)
	db.SetMaxIdleConns(1)
	db.SetConnMaxLifetime(0)

	// Ping the database to verify the connection
	for i := 0; i < 10; i++ {
		err = db.Ping()
		if err == nil {
			log.Println("Connected to the database.")
			return db, nil
		}
		log.Printf("Unable to connect to the database, retrying... (%d/10)", i+1)
		time.Sleep(1 * time.Second)
	}

	return nil, err
}
