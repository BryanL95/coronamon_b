package server

import (
	"fmt"
	"log"
	"net/http"

	"github.com/coronamon/data"
)

func Server(w http.ResponseWriter, r *http.Request) {
	_, err := data.LoadData()

	if err != nil {
		log.Println("Error to call function loadData")
	}

	fmt.Fprintf(w, "Response")
}
