package server

import (
	"log"
	"net/http"

	"github.com/coronamon/data"
)

func Server(w http.ResponseWriter, r *http.Request) {
	response, error := data.LoadData()

	if error != nil {
		log.Println("Error to response server, check loadData function")
		http.Error(w, error.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-type", "application/json")
	w.Write(response)
}
