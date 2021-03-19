package main

import (
	"database/sql"
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
	password = "palakopowe"
	host = "35.157.134.249"
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
	if _, err := db.Exec(createDatabase);
	err != nil {
		log.Fatal(err)
	}

	if _, err := db.Exec(tableCreationQuery);
	err != nil {
		log.Fatal(err)
	}
}

func CreateItem(w http.ResponseWriter, r *http.Request) {
    description := r.FormValue("description")
    log.WithFields(log.Fields{"description": description}).Info("Add new TodoItem. Saving to database.")
    todo := &TodoItemModel{Description: description, Completed: false}
	var lastInsertId int
	var err error
	err = db.QueryRow("INSERT INTO todolist(description, completed) VALUES($2, $3) returning id", todo).Scan(&lastInsertId)
	if err != nil {
		log.Fatal("Failed to execute query: ", err)
	}
	userSql := "SELECT * FROM todolist"
	err = db.QueryRow(userSql).Scan(&todo.Id, &todo.Description, &todo.Completed)
    w.Header().Set("Content-Type", "application/json")
}


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
	http.ListenAndServe(":8000", router)

}