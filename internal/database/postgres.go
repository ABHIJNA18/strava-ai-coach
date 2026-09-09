package database

import (
	"database/sql"
	"fmt"
	"os"

	_ "github.com/lib/pq"

)
// if DATABASE_URL is set, it will be used as the connection string. 
// Otherwise, the connection string will be constructed from individual environment variables.

func NewPostgresConnection() (*sql.DB, error) {
	databaseURL := os.Getenv("DATABASE_URL")

	var connectionString string

	if databaseURL != "" {
		connectionString = databaseURL
	} else {
		host := os.Getenv("DB_HOST")
		port := os.Getenv("DB_PORT")
		user := os.Getenv("DB_USER")
		password := os.Getenv("DB_PASSWORD")
		dbName := os.Getenv("DB_NAME")

		connectionString = fmt.Sprintf(
			"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
			host,
			port,
			user,
			password,
			dbName,
		)
	}

	db, err := sql.Open(
		"postgres",
		connectionString,
	)
	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		return nil, err
	}

	fmt.Println("Connection to database successful")

	return db, nil
}