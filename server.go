package main

import (
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
)

type Page struct {
	Title string
	Body  []byte
}

var templates = template.Must(template.ParseFiles("homepage.html", "error.html", "resources.html"))

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
	log.Fatal(http.ListenAndServe(":8080", nil))
}
