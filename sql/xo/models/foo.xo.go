// Package models contains the types for schema ''.
package models

// GENERATED BY XO. DO NOT EDIT.

import (
	"database/sql"
	"errors"
)

// Foo represents a row from 'foo'.
type Foo struct {
	ID   int            `json:"id"`   // id
	Name sql.NullString `json:"name"` // name

	// xo fields
	_exists, _deleted bool
}

// Exists determines if the Foo exists in the database.
func (f *Foo) Exists() bool {
	return f._exists
}

// Deleted provides information if the Foo has been deleted from the database.
func (f *Foo) Deleted() bool {
	return f._deleted
}

// Insert inserts the Foo to the database.
func (f *Foo) Insert(db XODB) error {
	var err error

	// if already exist, bail
	if f._exists {
		return errors.New("insert failed: already exists")
	}

	// sql query
	const sqlstr = `INSERT INTO foo (` +
		`name` +
		`) VALUES (` +
		`?` +
		`)`

	// run query
	XOLog(sqlstr, f.Name)
	res, err := db.Exec(sqlstr, f.Name)
	if err != nil {
		return err
	}

	// retrieve id
	id, err := res.LastInsertId()
	if err != nil {
		return err
	}

	// set primary key and existence
	f.ID = int(id)
	f._exists = true

	return nil
}

// Update updates the Foo in the database.
func (f *Foo) Update(db XODB) error {
	var err error

	// if doesn't exist, bail
	if !f._exists {
		return errors.New("update failed: does not exist")
	}

	// if deleted, bail
	if f._deleted {
		return errors.New("update failed: marked for deletion")
	}

	// sql query
	const sqlstr = `UPDATE foo SET ` +
		`name = ?` +
		` WHERE id = ?`

	// run query
	XOLog(sqlstr, f.Name, f.ID)
	_, err = db.Exec(sqlstr, f.Name, f.ID)
	return err
}

// Save saves the Foo to the database.
func (f *Foo) Save(db XODB) error {
	if f.Exists() {
		return f.Update(db)
	}

	return f.Insert(db)
}

// Delete deletes the Foo from the database.
func (f *Foo) Delete(db XODB) error {
	var err error

	// if doesn't exist, bail
	if !f._exists {
		return nil
	}

	// if deleted, bail
	if f._deleted {
		return nil
	}

	// sql query
	const sqlstr = `DELETE FROM foo WHERE id = ?`

	// run query
	XOLog(sqlstr, f.ID)
	_, err = db.Exec(sqlstr, f.ID)
	if err != nil {
		return err
	}

	// set deleted
	f._deleted = true

	return nil
}

// FooByID retrieves a row from 'foo' as a Foo.
//
// Generated from index 'foo_id_pkey'.
func FooByID(db XODB, id int) (*Foo, error) {
	var err error

	// sql query
	const sqlstr = `SELECT ` +
		`id, name ` +
		`FROM foo ` +
		`WHERE id = ?`

	// run query
	XOLog(sqlstr, id)
	f := Foo{
		_exists: true,
	}

	err = db.QueryRow(sqlstr, id).Scan(&f.ID, &f.Name)
	if err != nil {
		return nil, err
	}

	return &f, nil
}