package server

import (
	"log"
	"net/http"

	"github.com/coronamon/data"
)

func Server(w http.ResponseWriter, r *http.Request) {

	route := r.URL.RequestURI()
	w.Header().Set("Content-type", "application/json")

	switch route {
	case "/":
		response, error := data.GeneralData()

		if error != nil {
			log.Println("Error to response server, check GeneralData function")
			http.Error(w, error.Error(), http.StatusInternalServerError)
			return
		}
		w.Write(response)
	case "/last-day":
		response, error := data.LoadData()

		if error != nil {
			log.Println("Error to response server, check loadData function")
			http.Error(w, error.Error(), http.StatusInternalServerError)
			return
		}
		w.Write(response)
	case "/last-day-c":
		w.Write([]byte{})
	case "/country":
		w.Write([]byte{})
	}

}
