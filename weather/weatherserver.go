package main

import (
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
)

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

type Page struct {
	Title string
	Body  []byte
}

var templates = template.Must(template.ParseFiles("weather.html", "C:/Users/Aleksander/go/src/AgriMaz-ERP/weather/error.html" ))

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
	//getWeather("xxx", "xxx","hourly,daily,alerts")
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
	title := r.URL.Path[len("/"):]
	p, err := loadPage(title)
	if err != nil {
		http.Redirect(w, r, "/error", http.StatusFound)
		return
	}
	renderTemplate(w, title, p)
}

func main() {
	http.HandleFunc("/weather", mainHandler)
	http.HandleFunc("/error", errorHandler)
	log.Fatal(http.ListenAndServe(":8080", nil))
}