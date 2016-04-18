package main

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
)

const baseUrl = `https://xkcd.com/%d/info.0.json`

type ComicInfo struct {
	Num int
	Day,
	Year,
	Month,
	Link,
	News,
	Transcript,
	Alt,
	Img,
	Title string
	SafeTitle string `json:"safe_title"`
}

type FileInfo struct {
	JSONFilename,
	ImageFilename string
}

func main() {

	// Find the higest number we have completed and start there
	comicNumber := 1
	files, err := ioutil.ReadDir(".")
	if err != nil {
		fmt.Println("Error reading dir:", err)
		os.Exit(1)
	}
	for _, file := range files {
		if filepath.Ext(file.Name()) == ".jpg" || filepath.Ext(file.Name()) == ".png" {
			var i int
			fmt.Sscanf(file.Name(), "%d", &i)
			if i > comicNumber {
				comicNumber = i
			}
		}
	}

	// Fetch the images and json metadata
	var errCount int
	for i := comicNumber; ; i++ {
		comic, files, err := getComic(i)
		fmt.Printf("%+v\n%+v\n", comic, files)
		if err != nil {
			errCount++
			fmt.Println(err)
			if errCount > 2 {
				break
			}
		}
	}

}

func getComic(comicNumber int) (*ComicInfo, *FileInfo, error) {

	fileInfo := &FileInfo{}

	// Get the JSON
	q := fmt.Sprintf(baseUrl, comicNumber)
	fmt.Println("Getting:", q)
	resp, err := http.Get(q)
	if err != nil {
		return nil, nil, err
	}
	// We must close resp.Body on all execution paths.
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, nil, fmt.Errorf("query failed: %s", resp.Status)
	}

	// Store the JSON in a file as well
	infoFile, err := os.Create(fmt.Sprintf("./%d.json", comicNumber))
	if err != nil {
		return nil, nil, err
	}
	fileInfo.JSONFilename = fmt.Sprintf("./%d.json", comicNumber)

	// Write to file when this reader is read
	infoTeeReader := io.TeeReader(resp.Body, infoFile)

	var result ComicInfo
	if err := json.NewDecoder(infoTeeReader).Decode(&result); err != nil {
		return nil, nil, err
	}
	fmt.Printf("%+v\n\n", result)

	// Fetch the image
	if result.Img != "" && result.Img != "http://imgs.xkcd.com/comics/" {
		fmt.Println("Getting:", result.Img)
		imgRes, err := http.Get(result.Img)
		if err != nil {
			return nil, nil, err
		}
		// We must close resp.Body on all execution paths.
		defer imgRes.Body.Close()

		if imgRes.StatusCode != http.StatusOK {
			return nil, nil, fmt.Errorf("query failed: %s", imgRes.Status)
		}

		imageFilename := fmt.Sprintf("./%d%s", comicNumber, filepath.Ext(result.Img))
		imageFile, err := os.Create(imageFilename)

		if err != nil {
			return nil, nil, err
		}

		_, err = io.Copy(imageFile, imgRes.Body)
		if err != nil {
			return nil, nil, err
		}
		fileInfo.ImageFilename = imageFilename

	} else {
		fmt.Println("Not fetching:", result.Img)
	}

	return &result, fileInfo, nil
}
