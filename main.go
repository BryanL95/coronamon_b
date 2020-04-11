package main

import (
	"net/http"

	"github.com/coronamon/server"
)

//https://github.com/CSSEGISandData/COVID-19/tree/master/csse_covid_19_data/csse_covid_19_time_series --> Global by contry
//https://github.com/CSSEGISandData/COVID-19/tree/master/csse_covid_19_data/csse_covid_19_daily_reports --> Global by last day/country

func main() {
	http.HandleFunc("/", server.Server)
	http.ListenAndServe(":3000", nil)
}
