package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/drive/v3"
)

func main() {

	if len(os.Args) < 2 {
		fmt.Fprintln(os.Stderr, "Usage: drive filename (to upload a file)")
		return
	}
	ctx := context.Background()
	//ctx = oauth2.NoContext

	//filename := os.Args[1]

	client1, err := google.DefaultClient(ctx, drive.DriveScope)
	if err != nil {
		log.Fatalln(err)
	}

	ts, err := google.DefaultTokenSource(ctx, drive.DriveScope)
	if err != nil {
		log.Fatalln(err)
	}
	client2 := oauth2.NewClient(ctx, ts)
	_ = client2
	_ = client1

	//fmt.Printf("Client: %+v", client1.)

	service, err := drive.New(client2)
	if err != nil {
		log.Fatalf("Unable to create Drive service: %v", err)
	}

	a, e := service.About.Get().Do()
	fmt.Printf("About:\n%+v\n\nerr:\n%v\n\n", a, e)

	//goFile, err := os.Open(filename)
	//if err != nil {
	//	log.Fatalf("error opening %q: %v", filename, err)
	//}

	//fileService := drive.NewFilesService(service)
	//file := &drive.File{}
	//file.Name = filename
	//create := fileService.Create(file)
	//file2, err := create.Do()
	//log.Printf("err: %s, file: %+v", err, file2)
	//driveFile, err := service.Files.Insert(&drive.File{Title: filename}).Media(goFile).Do()

	//log.Printf("Got drive.File, err: %#v, %v", driveFile, err)

}
