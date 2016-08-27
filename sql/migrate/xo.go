package main

import (
	"database/sql"
	"log"

	_ "github.com/mattn/go-sqlite3"
	"github.com/oyvindsk/go-playground/sql/migrate/models"
)

func main() {
	db, err := sql.Open("sqlite3", "../foo2.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// create a post
	p := models.Post{Title: "Hello World", Body: "This is super duper exiting", CreatedAt: 0, UpdatedAt: 0}
	p.Save(db)

	// print all posts

}
