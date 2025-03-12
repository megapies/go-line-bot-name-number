package handler

import (
	"fmt"
	"log"
	"net/http"
)

func Handler(w http.ResponseWriter, r *http.Request) {
	log.Println("Request received", r.Method, r.URL.Path)
	fmt.Fprintf(w, "<h1>Hello from Go!</h1>")
}
