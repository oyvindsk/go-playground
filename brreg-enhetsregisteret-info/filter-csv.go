package main

// http://www.brreg.no/oppslag/
// http://data.brreg.no/oppslag/enhetsregisteret/enheter.xhtml
// https://confluence.brreg.no/display/DBNPUB/API
// https://confluence.brreg.no/display/DBNPUB/API#API-Filtrering

import (
	"fmt"
     "encoding/csv"
     "os"
)

func main() {

    // open the csv file
    f, err := os.Open("test.csv")
    if err != nil {
        fmt.Println("Could not open file: ", err)
    }

    // parse and filter data
    reader := csv.NewReader(f)
    reader.Comma = ';'
    for {
        record, err := reader.Read()
        if err != nil {
            fmt.Println("Could not parse line", err)
            break
        }
        fmt.Printf("%+v\n", record[0])
    }

}
