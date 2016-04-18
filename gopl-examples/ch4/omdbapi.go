// FIXME - Quote the url before getting the JSON

package main

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path"
	"path/filepath"
)

const baseApiUrl = "http://www.omdbapi.com/?t=%s&y=&plot=short&r=json"

func main() {

	if len(os.Args) < 2 {
		fmt.Println("Usage:", filepath.Base(os.Args[0]), "movie-title")
		os.Exit(2)
	}

	title := os.Args[1]
	fmt.Println(title)

	// Get the JSON
	q := fmt.Sprintf(baseApiUrl, title)
	fmt.Println("Getting:", q)
	resp, err := http.Get(q)
	if err != nil {
		log.Fatal("Could not fetch movie data:", err)
	}
	// We must close resp.Body on all execution paths.
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Fatal("query failed: %s", resp.Status)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}
	var movie struct{ Title, IMDBID, Poster string }
	if err := json.Unmarshal(body, &movie); err != nil {
		log.Fatal("Json unmarshal failed:", err)
	}

	posterFilename := fmt.Sprintf("%s-%s%s", movie.Title, movie.IMDBID, path.Ext(movie.Poster))
	fmt.Println("Getting:", movie.Poster, "===>", posterFilename)
	resp, err = http.Get(movie.Poster)
	if err != nil {
		log.Fatal("Could not fetch poster:", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		log.Fatal("query failed: %s", resp.Status)
	}

	posterFile, err := os.Create(posterFilename)
	if err != nil {
		log.Fatal("Could not create file", posterFilename, ":", err)
	}

	_, err = io.Copy(posterFile, resp.Body)
	if err != nil {
		log.Fatal("Could not write poster to file", posterFilename, ":", err)
	}

}
