package main

import (
	"context"
	"database/sql"
	"fmt"
	"io"
	"time"

	"github.com/labstack/echo/v4"
	_ "github.com/mattn/go-sqlite3"
)

func mustInitDB(logger echo.Logger, dbfile string) *sql.DB {

	db, err := sql.Open("sqlite3", dbfile)
	if err != nil {
		logger.Fatal(err)
	}

	sqlStmt := `CREATE TABLE IF NOT EXISTS messages (unixmilli INTEGER NO NULL PRIMARY KEY, author, msg text)`

	_, err = db.Exec(sqlStmt)
	if err != nil {
		logger.Fatalf("%q: %s\n", err, sqlStmt)
	}

	return db
}

type dbMessages []struct {
	time        int
	author, msg string
}

func getAll(db *sql.DB, logger echo.Logger) (dbMessages, error) {

	var res dbMessages

	rows, err := db.Query("SELECT unixmilli, author, msg FROM messages ORDER BY unixmilli DESC")
	if err != nil {
		return dbMessages{}, err
	}

	defer rows.Close()

	for rows.Next() {
		var time int
		var author string
		var msg string
		err = rows.Scan(&time, &author, &msg)
		if err != nil {
			return dbMessages{}, err
		}

		res = append(res, struct {
			time        int
			author, msg string
		}{time, author, msg})

	}

	err = rows.Err()
	if err != nil {
		return dbMessages{}, err
	}

	return res, nil
}

func fprintAll(logger echo.Logger, db *sql.DB, out io.Writer) error {

	rows, err := db.Query("SELECT unixmilli, author, msg FROM messages")
	if err != nil {
		return err
	}

	defer rows.Close()

	for rows.Next() {
		var time int
		var author string
		var msg string
		err = rows.Scan(&time, &author, &msg)
		if err != nil {
			return err
		}
		fmt.Fprintln(out, time, author, msg)
	}

	err = rows.Err()
	if err != nil {
		return err
	}

	return nil
}

func storeMsg(logger echo.Logger, db *sql.DB, author, msg string) error {

	ctx := context.Background()

	res, err := db.ExecContext(ctx, "INSERT INTO messages(unixmilli, author, msg) VALUES(?, ?, ?)", time.Now().UnixMilli(), author, msg)
	if err != nil {
		return err
	}

	rows, err := res.RowsAffected()

	if err != nil {
		logger.Infof("storeMsg: res: Error in RowsAffected: %s", err)
		return err
	}

	logger.Infof("storeMsg: res: %d", rows)

	return nil
}

func foo(srv server, logger echo.Logger) {

	// os.Remove("./foo.db")

	db := srv.db // =/

	tx, err := db.Begin()
	if err != nil {
		logger.Fatal(err)
	}
	stmt, err := tx.Prepare("insert into foo(id, name) values(?, ?)")
	if err != nil {
		logger.Fatal(err)
	}
	defer stmt.Close()
	for i := 0; i < 100; i++ {
		_, err = stmt.Exec(i, fmt.Sprintf("こんにちわ世界%03d", i))
		if err != nil {
			logger.Fatal(err)
		}
	}
	err = tx.Commit()
	if err != nil {
		logger.Fatal(err)
	}

	rows, err := db.Query("select id, name from foo")
	if err != nil {
		logger.Fatal(err)
	}
	defer rows.Close()
	for rows.Next() {
		var id int
		var name string
		err = rows.Scan(&id, &name)
		if err != nil {
			logger.Fatal(err)
		}
		fmt.Println(id, name)
	}
	err = rows.Err()
	if err != nil {
		logger.Fatal(err)
	}

	stmt, err = db.Prepare("select name from foo where id = ?")
	if err != nil {
		logger.Fatal(err)
	}
	defer stmt.Close()
	var name string
	err = stmt.QueryRow("3").Scan(&name)
	if err != nil {
		logger.Fatal(err)
	}
	fmt.Println(name)

	_, err = db.Exec("delete from foo")
	if err != nil {
		logger.Fatal(err)
	}

	_, err = db.Exec("insert into foo(id, name) values(1, 'foo'), (2, 'bar'), (3, 'baz')")
	if err != nil {
		logger.Fatal(err)
	}

	rows, err = db.Query("select id, name from foo")
	if err != nil {
		logger.Fatal(err)
	}
	defer rows.Close()
	for rows.Next() {
		var id int
		var name string
		err = rows.Scan(&id, &name)
		if err != nil {
			logger.Fatal(err)
		}
		fmt.Println(id, name)
	}
	err = rows.Err()
	if err != nil {
		logger.Fatal(err)
	}

}
