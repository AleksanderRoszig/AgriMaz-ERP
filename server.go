package main

import (
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
)

type Page struct {
	Title string
	Body  []byte
}

var templates = template.Must(template.ParseFiles("homepage.html", "error.html", "resources.html", "calendar.html", "weather.html"))

func getWeather(latitude string, longitude string, part string) {
	var url string
	url = "https://api.openweathermap.org/data/2.5/onecall?lat={lat}&lon={lon}&exclude={part}&appid={API key}"
	url = strings.ReplaceAll(url, "{lat}", latitude)
	url = strings.ReplaceAll(url, "{lon}", longitude)
	url = strings.ReplaceAll(url, "{part}", part)
	url = strings.ReplaceAll(url, "{API key}", "bb7f4bb3c58c951bb1847c8d0c222746")
	fmt.Println(url)
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
	getWeather("33.441792", "-94.037689","hourly,daily,alerts")
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
	log.Fatal(http.ListenAndServe(":8080", nil))
}