package todo

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
	log "github.com/sirupsen/logrus"
	"net/http"
	"strconv"
)
var db *sql.DB
const(
	Host = "18.198.177.82"
	Port = 5432
	User= "myuser"
	Password = "test"
	Dbname = "todo"
)
type App struct {
	Router *mux.Router
	DB     *sql.DB
}

func (a *App) Initialize(host string, port int, user string, password string, dbname string) {
	connectionString := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", Host, Port, User, Password, Dbname)

	var err error
	a.DB, err = sql.Open("postgres", connectionString)
	if err != nil {
		log.Fatal(err)
	}
	err = a.DB.Ping()
	if err != nil {
		panic(err)
	}
	a.Router = mux.NewRouter()

	a.initializeRoutes()
	fmt.Println("Successfully connected!")
}


func (a *App) Run(addr string) {
	log.Fatal(http.ListenAndServe(":8010", a.Router))
}

func (a *App) getTask(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid task ID")
		return
	}

	p := todoItemModel{Id: id}
	if err := p.getTask(a.DB); err != nil {
		switch err {
		case sql.ErrNoRows:
			respondWithError(w, http.StatusNotFound, "Task not found")
		default:
			respondWithError(w, http.StatusInternalServerError, err.Error())
		}
		return
	}

	respondWithJSON(w, http.StatusOK, p)
}


func respondWithError(w http.ResponseWriter, code int, message string) {
	respondWithJSON(w, code, map[string]string{"error": message})
}

func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	response, _ := json.Marshal(payload)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
}

func (a *App) getTasks(w http.ResponseWriter, r *http.Request) {
	count, _ := strconv.Atoi(r.FormValue("count"))
	start, _ := strconv.Atoi(r.FormValue("start"))

	if count > 10 || count < 1 {
		count = 10
	}
	if start < 0 {
		start = 0
	}

	products, err := getTasks(a.DB, start, count)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, products)
}

func (a *App) createTask(w http.ResponseWriter, r *http.Request) {
	var p todoItemModel
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&p); err != nil {
		fmt.Println("cos poszlo xle")
		respondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}
	defer r.Body.Close()

	if err := p.createTask(a.DB); err != nil {
		fmt.Println("rerrorr")
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusCreated, p)
}


func (a *App) updateTask(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		fmt.Println("blad invalid task ID")
		respondWithError(w, http.StatusBadRequest, "Invalid task ID")
		return
	}

	var p todoItemModel
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&p); err != nil {
		fmt.Println("blad invalid task ID2")
		respondWithError(w, http.StatusBadRequest, "Invalid resquest payload")
		return
	}
	defer r.Body.Close()
	p.Id = id

	if err := p.updateTask(a.DB); err != nil {
		fmt.Println("blad invalid task ID3")
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, p)
}

func (a *App) deleteTask(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid Task ID")
		return
	}

	p := todoItemModel{Id: id}
	if err := p.deleteTask(a.DB); err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, map[string]string{"result": "success"})
}


func (a *App) initializeRoutes() {
	a.Router.HandleFunc("/tasks", a.getTasks).Methods("GET")
	a.Router.HandleFunc("/task", a.createTask).Methods("POST")
	a.Router.HandleFunc("/task/{id:[0-9]+}", a.getTask).Methods("GET")
	a.Router.HandleFunc("/task/{id:[0-9]+}", a.updateTask).Methods("PUT")
	a.Router.HandleFunc("/task/{id:[0-9]+}", a.deleteTask).Methods("DELETE")
}
func init() {
	log.SetFormatter(&log.TextFormatter{})
	log.SetReportCaller(true)
}

func main() {
	log.Info("Starting Todolist API server")
	a := App{}
	a.Initialize(Host, Port, User, Password, Dbname)

	a.Run(":8010")

}