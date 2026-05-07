package database

import (
	"log"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

func Connect(connURL string) (*sqlx.DB, error) {
	db, err := sqlx.Connect("postgres", connURL)
	if err != nil {
		log.Printf("Failed to connect to database: %v", err)
		return nil, err
	}

	if err := db.Ping(); err != nil {
		log.Printf("Failed to ping database: %v", err)
		return nil, err
	}

	log.Println("Successfully connected to the database")
	return db, nil
}
