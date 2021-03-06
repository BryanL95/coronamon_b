package main

import (
	"net/http"

	"github.com/coronamon/server"
	"github.com/gorilla/mux"
)

//https://github.com/CSSEGISandData/COVID-19/tree/master/csse_covid_19_data/csse_covid_19_time_series --> Global by contry
//https://github.com/CSSEGISandData/COVID-19/tree/master/csse_covid_19_data/csse_covid_19_daily_reports --> Global by last day/country

func main() {
	r := mux.NewRouter()

	r.HandleFunc("/", server.Global)                        //Load global data - day by day
	r.HandleFunc("/global/{country}", server.Global)        //Load data by last day confirmed general
	r.HandleFunc("/last-day/{date}", server.Last)           //Load data by last day confirmed general
	r.HandleFunc("/last-day/{date}/{country}", server.Last) //load data by country from last day

	http.ListenAndServe(":3001", r) //Pass handler (r) to server
	/*err := http.ListenAndServe(getPort(), r) //Pass handler (r) to server
	if err != nil {
		log.Println("ListenAndServe", err)
	}*/
}

/*
* Only enabled for heroku
 */
/*func getPort() string {
	var port = os.Getenv("PORT")
	if port == "" {
		port = "4747"
	}
	return ":" + port
}*/
