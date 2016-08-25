package middlewareTest

import (
	"context"
	"log"
	"net/http"
)

type Middleware struct {
	isFoo bool
	data  string
}

// Middleware is a struct that has a ServeHTTP method
func NewMiddleware() *Middleware {
	return &Middleware{isFoo: true}
}

// The middleware handler
func (m *Middleware) ServeHTTP(w http.ResponseWriter, req *http.Request, next http.HandlerFunc) {
	// Log moose status
	log.Printf("Middleware BEFORE: %v %v\n", m.isFoo, m.data)

	w.Write([]byte("Hello from MiddleWare on the way in\n"))

	// Use context values only for request-scoped data that transits
	// processes and API boundaries, not for passing optional parameters to
	// functions.
	// https://godoc.org/context#Context
	ctx := context.WithValue(req.Context(), "foo", "BAR")
	req = req.WithContext(ctx)

	// Call the next middleware handler
	next(w, req)

	w.Write([]byte("Hello MiddleWare on the way out\n"))

	log.Printf("Middleware AFTER: %v %v\n", m.isFoo, m.data)

}
