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
	//"database/sql"
)

const tableCreationQuery = `CREATE TABLE IF NOT EXISTS todolist
(
    
    name TEXT NOT NULL,
    price NUMERIC(10,2) NOT NULL DEFAULT 0.00,
    CONSTRAINT products_pkey PRIMARY KEY (id)
)`

type TodoItemModel struct{
       Id int
       Description string
       Completed bool
}

func CreateItem(w http.ResponseWriter, r *http.Request) {
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


    description := r.FormValue("description")
    log.WithFields(log.Fields{"description": description}).Info("Add new TodoItem. Saving to database.")
    todo := &TodoItemModel{Description: description, Completed: false}
	var lastInsertId int
	err = db.QueryRow("INSERT INTO todolist(name,date,price) VALUES($1,$2) returning id", "car", "19.03.2021", "2332").Scan(&lastInsertId)
	if err != nil {
		log.Fatal("Failed to execute query: ", err)
	}
	userSql := "SELECT * FROM todolist"
	err = db.QueryRow(userSql).Scan(&todo.Id, &todo.Description, &todo.Completed)

    //result := db.Last(&todo)
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(result.Value)
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

func main() {
	log.Info("Starting Todolist API server")
	router := mux.NewRouter()
	router.HandleFunc("/healthz", Healthz).Methods("GET")
	http.ListenAndServe(":8000", router)
}