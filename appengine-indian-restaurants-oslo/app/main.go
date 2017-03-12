package app

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/gorilla/mux"
	"google.golang.org/appengine"
	"google.golang.org/appengine/log"
)

// Init, run by AE
func init() {
	r := mux.NewRouter()
	r.Handle("/", makeHandler(handlePageMain))

	http.Handle("/", r)
}

// Page Handler
func handlePageMain(w http.ResponseWriter, r *http.Request) error {
	ctx := appengine.NewContext(r)

	visited, unvisited, err := getRestaurantsByVisited(ctx)

	if err != nil {
		return err
	}

	fmt.Fprintf(w, "Visted:\n%s\n\nNot Visited yet:\n%s\n", strings.Join(visited, "\n"), strings.Join(unvisited, "\n"))
	return nil
}

// Helper functions
func makeHandler(fn func(http.ResponseWriter, *http.Request) error) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := appengine.NewContext(r)
		err := fn(w, r)

		if err != nil {
			if ferr, ok := err.(fooError); ok {
				log.Errorf(ctx, "!! %d: %s (%v)", ferr.HTTPCode, ferr.Msg, ferr.Origin)
				http.Error(w, ferr.Msg, ferr.HTTPCode)
				// could this be moved to middleware?
			} else {
				log.Errorf(ctx, "!! Error: %s", err)
			}
		}
	}
}

// Custom errror type
type fooError struct {
	Origin   error
	Msg      string
	HTTPCode int
}

func (fe fooError) Error() string {
	return fmt.Sprintf("%d: %s (%s)", fe.HTTPCode, fe.Msg, fe.Origin.Error())
}
