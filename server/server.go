package server

import (
	"log"
	"net/http"

	"github.com/coronamon/data"
	"github.com/gorilla/mux"
)

func Global(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Content-type", "application/json")
	params := mux.Vars(r)

	if len(params) < 1 {
		response, error := data.GeneralData("")

		if error != nil {
			log.Println("Error to response server, check GeneralData function")
			http.Error(w, error.Error(), http.StatusInternalServerError)
			return
		}
		w.Write(response)
	} else {
		response, error := data.GeneralData(params["country"])

		if error != nil {
			log.Println("Error to response server, check GeneralData function")
			http.Error(w, error.Error(), http.StatusInternalServerError)
			return
		}
		w.Write(response)
	}
}

func Last(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Content-type", "application/json")
	params := mux.Vars(r)

	response, error := data.LoadData(params["country"], params["date"])

	if error != nil {
		log.Println("Error to response server, check loadData function")
		http.Error(w, error.Error(), http.StatusInternalServerError)
		return
	}
	w.Write(response)
}
