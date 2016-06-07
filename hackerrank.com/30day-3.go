package main

import (
	"fmt"
	"log"
)

func main() {
	var n int
	_, err := fmt.Scan(&n)
	if err != nil {
		log.Fatal(err)
	}

	switch {
	case n%2 != 0, n >= 6 && n <= 20:
		fmt.Println("Weird")
	case n >= 2 && n <= 5, n > 20:
		fmt.Println("Not Weird")
	}

}
