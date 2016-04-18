// Issues prints a table of GitHub issues matching the search terms.
package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/oyvindsk/go-playground/gopl-examples/ch4/github"
)

func main() {
	result, err := github.SearchIssues(os.Args[1:])
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%d issues:\n", result.TotalCount)
	now := time.Now()
	for _, item := range result.Items {
		age := now.Sub(item.UpdatedAt)
		var ageS string
		switch {
		case age.Hours() < 24*30:
			ageS = "Less than a month old"
		case age.Hours() < 24*30*365:
			ageS = "Less than a year old"
		default:
			ageS = "More than a year old"
		}

		fmt.Printf("#%-5d %9.9s %.55s\t%v\t%v\n", item.Number, item.User.Login, item.Title, item.UpdatedAt, ageS)
	}
}
