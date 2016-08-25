package main

// https://medium.com/@matryer/context-has-arrived-per-request-state-in-go-1-7-4d095be83bd8
// https://www.reddit.com/r/golang/comments/4tm847/how_to_correctly_use_contextcontext_in_go_17/
// https://www.reddit.com/r/golang/comments/4s1954/context_has_arrived_perrequest_state_in_go_17/
// https://godoc.org/context#Context
// https://elithrar.github.io/article/map-string-interface/

import (
	"fmt"
	"log"
	"net/http"

	"github.com/oyvindsk/go-playground/negroni-context/middlewareTest"
	"github.com/urfave/negroni"
)

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, req *http.Request) {
		fmt.Fprintf(w, "Welcome to the home page!\n")

		var valueFromMiddleware string
		valueFromMiddleware, ok := req.Context().Value("foo").(string)
		if ok {
			fmt.Fprintf(w, "got middleware value: %v\n", valueFromMiddleware)
		}
		log.Printf("main page!")
	})

	n := negroni.Classic() // Includes some default middlewares

	n.Use(middlewareTest.NewMiddleware())

	n.UseHandler(mux)

	http.ListenAndServe(":3000", n)

}
