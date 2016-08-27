// Package models contains the types for schema ''.
package models

// GENERATED BY XO. DO NOT EDIT.

import "errors"

// Author represents a row from 'authors'.
type Author struct {
	AuthorID int    `json:"author_id"` // author_id
	Name     string `json:"name"`      // name

	// xo fields
	_exists, _deleted bool
}

// Exists determines if the Author exists in the database.
func (a *Author) Exists() bool {
	return a._exists
}

// Deleted provides information if the Author has been deleted from the database.
func (a *Author) Deleted() bool {
	return a._deleted
}

// Insert inserts the Author to the database.
func (a *Author) Insert(db XODB) error {
	var err error

	// if already exist, bail
	if a._exists {
		return errors.New("insert failed: already exists")
	}

	// sql query
	const sqlstr = `INSERT INTO authors (` +
		`name` +
		`) VALUES (` +
		`?` +
		`)`

	// run query
	XOLog(sqlstr, a.Name)
	res, err := db.Exec(sqlstr, a.Name)
	if err != nil {
		return err
	}

	// retrieve id
	id, err := res.LastInsertId()
	if err != nil {
		return err
	}

	// set primary key and existence
	a.AuthorID = int(id)
	a._exists = true

	return nil
}

// Update updates the Author in the database.
func (a *Author) Update(db XODB) error {
	var err error

	// if doesn't exist, bail
	if !a._exists {
		return errors.New("update failed: does not exist")
	}

	// if deleted, bail
	if a._deleted {
		return errors.New("update failed: marked for deletion")
	}

	// sql query
	const sqlstr = `UPDATE authors SET ` +
		`name = ?` +
		` WHERE author_id = ?`

	// run query
	XOLog(sqlstr, a.Name, a.AuthorID)
	_, err = db.Exec(sqlstr, a.Name, a.AuthorID)
	return err
}

// Save saves the Author to the database.
func (a *Author) Save(db XODB) error {
	if a.Exists() {
		return a.Update(db)
	}

	return a.Insert(db)
}

// Delete deletes the Author from the database.
func (a *Author) Delete(db XODB) error {
	var err error

	// if doesn't exist, bail
	if !a._exists {
		return nil
	}

	// if deleted, bail
	if a._deleted {
		return nil
	}

	// sql query
	const sqlstr = `DELETE FROM authors WHERE author_id = ?`

	// run query
	XOLog(sqlstr, a.AuthorID)
	_, err = db.Exec(sqlstr, a.AuthorID)
	if err != nil {
		return err
	}

	// set deleted
	a._deleted = true

	return nil
}

// AuthorByAuthorID retrieves a row from 'authors' as a Author.
//
// Generated from index 'authors_author_id_pkey'.
func AuthorByAuthorID(db XODB, authorID int) (*Author, error) {
	var err error

	// sql query
	const sqlstr = `SELECT ` +
		`author_id, name ` +
		`FROM authors ` +
		`WHERE author_id = ?`

	// run query
	XOLog(sqlstr, authorID)
	a := Author{
		_exists: true,
	}

	err = db.QueryRow(sqlstr, authorID).Scan(&a.AuthorID, &a.Name)
	if err != nil {
		return nil, err
	}

	return &a, nil
}
