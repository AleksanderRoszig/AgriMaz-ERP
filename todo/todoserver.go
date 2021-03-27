package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
	log "github.com/sirupsen/logrus"
	"io"
	"net/http"
)
var db *sql.DB
const(
	user= "myuser"
	password = "test"
	host = "3.65.21.44"
	port = 5432
)
const createDatabase = `CREATE DATABASE tests;`

const tableCreationQuery = `CREATE TABLE IF NOT EXISTS todolist
(
    id SERIAL PRIMARY KEY, 
    description varchar(40) NOT NULL,
    completed integer NOT NULL
)`

type TodoItemModel struct{
       Id int
       Description string
       Completed bool
}
func ensureTableExists() {

	/*
	if _, err := db.Exec(createDatabase);
	err != nil {
		log.Fatal(err)
	}
*/
	if _, err := db.Exec(tableCreationQuery);
	err != nil {
		log.Fatal(err)
	}
}

func createItem(w http.ResponseWriter, r *http.Request) {
	var todo TodoItemModel
	err := json.NewDecoder(r.Body).Decode(&todo)
	if err != nil {
		http.Error(w,err.Error(),http.StatusBadRequest)
		return
	}
	fmt.Println(todo)
	var lastInsertId int

	err = db.QueryRow("INSERT INTO todolist(description, completed) VALUES($2, $3) returning id", todo).Scan(&lastInsertId)
	if err != nil {
		log.Fatal("Failed to execute query: ", err)
	}
	userSql := "SELECT * FROM todolist"
	err = db.QueryRow(userSql).Scan(&todo.Id, &todo.Description, &todo.Completed)
    w.Header().Set("Content-Type", "application/json")

    }
  //   */

func Healthz(w http.ResponseWriter, r *http.Request) {
	log.Info("API Health is OK")
	w.Header().Set("Content-Type", "application/json")
	io.WriteString(w, `{"alive": true}`)
}

func init() {
	log.SetFormatter(&log.TextFormatter{})
	log.SetReportCaller(true)
}
func connectToDB() {
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s sslmode=disable",
		host, port, user, password)
	var err error
	db, err = sql.Open("postgres", psqlInfo)
	if err != nil {
		panic(err)
	}
	err = db.Ping()
	if err != nil {
		panic(err)
	}

	fmt.Println("Successfully connected!")
}


func main() {
	log.Info("Starting Todolist API server")
	connectToDB()
    ensureTableExists()
	router := mux.NewRouter()
	router.HandleFunc("/healthz", Healthz).Methods("GET")
	router.HandleFunc("/todo", createItem).Methods("POST")
	http.ListenAndServe(":8000", router)

}