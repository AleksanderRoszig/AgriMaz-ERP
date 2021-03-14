package main

import (
	"database/sql"
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	_ "github.com/lib/pq"
)

type Page struct {
	Title string
	Body  []byte
}

type Purchase struct {
	Name string
	Date string
	Price float32
}

type weatherData struct {
	City      string
	Sunrise   int64
	Sunset    int64
	Temp      int
	FeelLike  int
	Pressure  int
	Humidity  int
	DewPoint  int
	Uvi       int
	Clouds    int
	WindSpeed int
	WindDeg   int
}

var templates = template.Must(template.ParseFiles("homepage.html", "error.html", "resources.html", "calendar.html", "weather.html", "purchases.html"))

func getFromDatabase(){
	const(
		user= "xxx"
		password = "xxx"
		host = "xxx"
		dbname = "xxx"
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
	var myThing Purchase
	userSql := "SELECT * FROM wydatki"
	err = db.QueryRow(userSql).Scan(&myThing.Name,&myThing.Date, &myThing.Price)
	if err != nil {
		log.Fatal("Failed to execute query: ", err)
	}
	fmt.Printf(myThing.Name)
}

func insertIntoDatabase(){
	const(
		user= "xxx"
		password = "xxx"
		host = "xxx"
		dbname = "xxx"
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
	var lastInsertId int
	err = db.QueryRow("INSERT INTO company(username,departname,created) VALUES($1,$2,$3) returning uid;", "car", "19.03.2021", "2332").Scan(&lastInsertId)
	if err != nil {
		log.Fatal("Failed to execute query: ", err)
		}

	fmt.Println("last inserted id =", lastInsertId)




}


func getJson(url string)(datafromURL string) {
	var bodyString string
	r, err := http.Get(url)
	if err != nil {
		log.Fatal(err)
	}

	defer r.Body.Close()

	if r.StatusCode == http.StatusOK {
		bodyBytes, err := ioutil.ReadAll(r.Body)
		if err != nil {
			log.Fatal(err)
		}
		bodyString = string(bodyBytes)
		fmt.Printf(bodyString)
	}

	return bodyString
}



func getWeather(latitude string, longitude string, part string) {
	var url string
	var jsonfromURL string
	url = "https://api.openweathermap.org/data/2.5/onecall?lat={lat}&lon={lon}&exclude={part}&appid={API key}"
	url = strings.ReplaceAll(url, "{lat}", latitude)
	url = strings.ReplaceAll(url, "{lon}", longitude)
	url = strings.ReplaceAll(url, "{part}", part)
	url = strings.ReplaceAll(url, "{API key}", "xxx")
	fmt.Println(url)
	jsonfromURL = getJson(url)
	fmt.Println(jsonfromURL)
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

func weatherHandler(w http.ResponseWriter, r *http.Request) {
	title := r.URL.Path[len("/"):]
	p, err := loadPage(title)
	if err != nil {
		http.Redirect(w, r, "/error", http.StatusFound)
		return
	}
	renderTemplate(w, title, p)
	getWeather("xxx", "xxx","hourly,daily,alerts")

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
	http.HandleFunc("/weather", weatherHandler)
	http.HandleFunc("/purchases", expensesHandler)
	log.Fatal(http.ListenAndServe(":8080", nil))
}