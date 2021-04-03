package weather

import (
	"database/sql"
	"fmt"
	"github.com/gorilla/mux"
	"io/ioutil"
	"net/http"
	"strings"
)

type App struct {
	Router *mux.Router
	DB     *sql.DB
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