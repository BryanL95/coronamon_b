package main

import (
	"net/http"

	"github.com/coronamon/server"
)

//https://github.com/CSSEGISandData/COVID-19/tree/master/csse_covid_19_data/csse_covid_19_time_series --> Global by contry
//https://github.com/CSSEGISandData/COVID-19/tree/master/csse_covid_19_data/csse_covid_19_daily_reports --> Global by last day/country

func main() {
	http.HandleFunc("/", server.Server)           //Load global data
	http.HandleFunc("/last-day", server.Server)   //Load data by last day
	http.HandleFunc("/last-day-c", server.Server) //load data by country from last day
	http.HandleFunc("/country", server.Server)    //Load data by country global
	http.ListenAndServe(":3000", nil)
}
