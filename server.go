package AgriMaz_ERP

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
)

type Page struct {
	Title string
	Body  []byte
}

type Purchase struct {
	id int
	Name string
	Date string
	Price float32
}

const(
	user= ""
	password = ""
	host = ""
	dbname = "company"
	port = 5432
)
var templates = template.Must(template.ParseFiles("homepage.html", "error.html", "resources.html", "calendar.html", "weather.html", "purchases.html"))

func getFromDatabase(){
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
	var myThing Purchase
	userSql := "SELECT * FROM expenses"
	err = db.QueryRow(userSql).Scan(&myThing.id, &myThing.Name, &myThing.Date, &myThing.Price)
	if err != nil {
		log.Fatal("Failed to execute query: ", err)
	}
	fmt.Printf(myThing.Name)
}

func insertIntoDatabase(){
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
	var lastInsertId int
	err = db.QueryRow("INSERT INTO expenses(name,date,price) VALUES($1,$2,$3) returning id", "car", "19.03.2021", "2332").Scan(&lastInsertId)
	if err != nil {
		log.Fatal("Failed to execute query: ", err)
	}
}



func (p *Page) save() error {
	filename := p.Title + ".txt"
	return ioutil.WriteFile(filename, p.Body, 0600)
}

func renderTemplate(w http.ResponseWriter, tmpl string, p *Page) {
	err := templates.ExecuteTemplate(w, tmpl+".html", p)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func loadPage(title string) (*Page, error) {
	filename := title + ".html"
	body, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	return &Page{Title: title, Body: body}, nil
}



func expensesHandler(w http.ResponseWriter, r *http.Request) {
	title := r.URL.Path[len("/"):]
	p, err := loadPage(title)
	if err != nil {
		http.Redirect(w, r, "/error", http.StatusFound)
		return
	}
	renderTemplate(w, title, p)
	getFromDatabase()
	//insertIntoDatabase()
}

func errorHandler(w http.ResponseWriter, r *http.Request) {
	title := r.URL.Path[len("/"):]
	p, err := loadPage(title)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	renderTemplate(w, title, p)
}

func mainHandler(w http.ResponseWriter, r *http.Request) {
	title := r.URL.Path[len("/mainpage/"):]
	p, err := loadPage(title)
	if err != nil {
		http.Redirect(w, r, "/error", http.StatusFound)
		return
	}
	renderTemplate(w, title, p)
}

func main() {
	http.HandleFunc("/", mainHandler)
	http.HandleFunc("/error", errorHandler)
	http.HandleFunc("/purchases", expensesHandler)
	log.Fatal(http.ListenAndServe(":8080", nil))
}