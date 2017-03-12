package app

import (
	"fmt"
	"net/http"
	"os"
	"strings"

	"golang.org/x/net/context"
	"google.golang.org/api/sheets/v4"
	"google.golang.org/appengine/log"
)

var sheetsClient *sheets.Service // ATM this is global and lazy initialzied by te functions in this file

func getRestaurantsByVisited(ctx context.Context) ([]string, []string, error) {

	if sheetsClient == nil {
		var err error
		sheetsClient, err = newSheetsClient(ctx)
		if err != nil {
			return nil, nil, err
		}
	}

	// Prints the names and majors of students in a sample spreadsheet:
	// https://docs.google.com/spreadsheets/d/1BxiMVs0XRA5nFMdKvBdBZjgmUUqptlbs74OgvE2upms/edit

	readRange := "Restauranter!A4:F"
	resp, err := sheetsClient.Spreadsheets.Values.Get(os.Getenv("SHEET_ID"), readRange).Do()
	if err != nil {
		log.Errorf(ctx, "Unable to retrieve data from sheet. %v", err)
		return nil, nil, fooError{Origin: err, Msg: "Unable to retrieve data from sheet", HTTPCode: http.StatusInternalServerError}
	}

	var visited []string
	var unvisited []string

	if len(resp.Values) > 0 {
		for _, row := range resp.Values {
			if len(row) == 0 {
				continue
			}

			if len(row) < 4 {
				log.Errorf(ctx, "! Could not parse row: %#v\n", row)
			}

			fmt.Printf("row:(%d)\n%#v\n", len(row), row)
			if row[0] != "" {
				s, ok := row[1].(string)
				if !ok {
					log.Errorf(ctx, "! Could not convert restaurant name to string: %v", row[1])
					continue
				}
				visited = append(visited, s)
			}
		}
	} else {
		fmt.Print("No data found.")
	}

	fmt.Printf("\n\nVisited:\n%s\n", strings.Join(visited, "\n"))

	return visited, unvisited, nil

}
