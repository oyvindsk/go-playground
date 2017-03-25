package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"golang.org/x/oauth2/google"
	"google.golang.org/api/compute/v1"
	"google.golang.org/api/drive/v3"
)

func main() {

	if len(os.Args) != 1 {
		fmt.Fprintln(os.Stderr, "Usage: drive filename (to upload a file)")
		return
	}
	ctx := context.Background()

	filename := os.Args[0]

	client, err := google.DefaultClient(ctx, compute.ComputeScope)
	if err != nil {
		log.Fatalln(err)
	}

	service, err := drive.New(client)
	if err != nil {
		log.Fatalf("Unable to create Drive service: %v", err)
	}

	goFile, err := os.Open(filename)
	if err != nil {
		log.Fatalf("error opening %q: %v", filename, err)
	}

	fileService := drive.NewFilesService(service)
	file := &drive.File{}
	file.Name = filename
	create := fileService.Create(file)
	file2, err := create.Do()
	log.Printf("err: %s", err)
	//driveFile, err := service.Files.Insert(&drive.File{Title: filename}).Media(goFile).Do()

	//log.Printf("Got drive.File, err: %#v, %v", driveFile, err)

}
