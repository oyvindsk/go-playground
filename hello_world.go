package main

import (
	"fmt"
	"log"
	"net/http"
)

func main() {

	http.HandleFunc(
		"/",
		func(w http.ResponseWriter, r *http.Request) {
			fmt.Fprintf(w, "<h1>Hello</h1>!")
		})

	log.Println("Listening..")
	http.ListenAndServe(":3000", nil)
}
