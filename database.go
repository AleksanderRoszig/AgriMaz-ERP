package main

import (
	"database/sql"
	"fmt"
	//"github.com/lib/pq"
)

const (
		hostname = "localhost"
		host_port = 5432
		username = "postgres"
		password = "test"
		database_name = "testdatabase"
	)

func main() {

	const(
		user= "test"
		password = "123456789"
		host = "localhost"
		dbname = "testdb"
		port = 5432
		)

	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)

	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	err = db.Ping()
	if err != nil {
		panic(err)
	}

	fmt.Println("Successfully connected!")
}