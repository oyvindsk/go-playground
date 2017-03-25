package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path"
)

func main() {

	// http://urort.p3.no/breeze/urort/TrackDTOViews?$filter=Recommended%20ne%20null&$orderby=Recommended%20desc%2CId%20desc&$top=24&$expand=Tags%2CFiles&$inlinecount=allpages&playlistId=0
	//var skip = 0

	resp, err := http.Get("http://urort.p3.no/breeze/urort/TrackDTOViews?$filter=Recommended%20ne%20null&$orderby=Recommended%20desc%2CId%20desc&$top=24&$expand=Tags%2CFiles&$inlinecount=allpages&playlistId=0&$skip=16")
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	dec := json.NewDecoder(resp.Body)
	for {
		var ur urørtResponse
		if err := dec.Decode(&ur); err == io.EOF {
			break
		} else if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("ID: %s, Typer: %s\n", ur.ID, ur.Type)

		for _, r := range ur.Results {
			fmt.Printf("IDs: %s, ID: %d, Title: %s\n", r.IDs, r.ID, r.Title)

			resp, err := http.Get(fmt.Sprintf("http://urort.p3.no/api/track/Download?trackId=%d", r.ID))
			if err != nil {
				log.Fatal(err)
			}
			defer resp.Body.Close()

			file, err := os.Create(path.Join("musikk", fmt.Sprintf("%s - %s.mp3", r.BandName, r.Title)))
			if err != nil {
				log.Fatal(err)
			}
			defer file.Close()

			io.Copy(file, resp.Body)
		}
	}
}
