package database

import (
	"database/sql"
	"fmt"
	"os"

	_ "github.com/lib/pq"

)

func NewPostgresConnection() (*sql.DB, error){

	host := os.Getenv("DB_HOST")
	port := os.Getenv("DB_PORT")
	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")
	dbName := os.Getenv("DB_NAME")

	connectionString := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		host,
		port,
		user,
		password,
		dbName,
	)

	//create a client to interact with the db
	db, err := sql.Open("postgres",connectionString)
	if err != nil{
		return nil, err
	}

	//interacct with db 
	err = db.Ping()
	if err != nil{
		return nil, err
	}

	fmt.Println("Connection to database successfull")

	return db, nil

}