package main

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/mattn/go-sqlite3"
	"github.com/oyvindsk/go-playground/sql/models"
)

func main() {
	db, err := sql.Open("sqlite3", "./foo.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	var foo *models.Foo
	foo, err = models.FooByID(db, 1)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%v : %s\n", foo.ID, foo.Name.String)

	// new foo
	foo2 := models.Foo{Name: sql.NullString{String: "new test", Valid: true}}
	err = foo2.Save(db)
	if err != nil {
		log.Fatal(err)
	}

	// new author
	ibsen, err := models.AuthorByAuthorID(db, 2)
	if err != nil {
		log.Fatal(err)
	}
	ibsen.Name = "Henrik Ibsen 2 2"
	ibsen.Save(db)

}
