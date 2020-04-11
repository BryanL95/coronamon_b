package main

import (
	"net/http"

	"github.com/coronamon/server"
)

func main() {
	http.HandleFunc("/", server.Server)
	http.ListenAndServe(":3000", nil)
}
