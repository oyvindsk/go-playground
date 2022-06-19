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

func mustInitDB(logger echo.Logger) *sql.DB {

	db, err := sql.Open("sqlite3", "./foo.db")
	if err != nil {
		logger.Fatal(err)
	}
	// defer db.Close()

	sqlStmt := `CREATE TABLE IF NOT EXISTS foo (unixmilli INTEGER NO NULL PRIMARY KEY, msg text)`

	// delete from fooo;

	_, err = db.Exec(sqlStmt)
	if err != nil {
		logger.Fatalf("%q: %s\n", err, sqlStmt)
	}

	return db
}

func fprintAll(logger echo.Logger, db *sql.DB, out io.Writer) error {

	rows, err := db.Query("SELECT unixmilli, msg FROM foo")
	if err != nil {
		return err
	}

	defer rows.Close()

	for rows.Next() {
		var id int
		var name string
		err = rows.Scan(&id, &name)
		if err != nil {
			return err
		}
		fmt.Fprintln(out, id, name)
	}

	err = rows.Err()
	if err != nil {
		return err
	}

	return nil
}

func storeMsg(logger echo.Logger, db *sql.DB, msg string) error {

	ctx := context.Background()

	res, err := db.ExecContext(ctx, "INSERT INTO foo(unixmilli, msg) VALUES(?, ?)", time.Now().UnixMilli(), msg)
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

func foo(logger echo.Logger) {

	// os.Remove("./foo.db")

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
