package main

import (
	"log"
)

// - Begynn på start (hviss error abort)
// - Let etter cyclecross linken  (hviss error abort)
// - Let etter riktig sykkel linken  (hviss error abort)
// - Velg størrelse  (hviss error abort)
// - Hent lagerstatus for alle butikker  (hviss error abort)
// - grep ut de i Oslo  (hviss error abort)
// - Mail output ??
//
//  Errors: Mail

func main() {
	err := findBikes()
	if err != nil {
		log.Println("\n\nERROR:\n", err, "\n!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!")
	}
}
