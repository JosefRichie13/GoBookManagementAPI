package main

import (
	"database/sql"

	_ "modernc.org/sqlite"

	"github.com/gin-gonic/gin"
)

// Defining JSON body for addNote() and addToANote(). It requires 2 JSON key's bookID, note text.
type AddNoteParameters struct {
	BookID string `json:"bookID" binding:"required"`
	Note   string `json:"note" binding:"required"`
}

// Adds a note to a Book, will clear all the old notes
func addNote(c *gin.Context) {

	// Variables for DB and Error
	var db *sql.DB
	var err error

	// Creating an instance of the struct, AddNoteParameters
	var addNoteParameters AddNoteParameters

	// Bind to the struct's members. If any member is invalid, binding does not happen and an error will be returned. Then its rejected with 400
	if c.BindJSON(&addNoteParameters) != nil {
		c.JSON(400, gin.H{"status": "Incorrect parameters, please provide all required parameters"})
		return
	}

	// Connect to the DB. If there is any issue connecting to the DB, throw a 500 error and return
	db, err = sql.Open("sqlite", "./BOOKMANAGEMENT.db")
	if err != nil {
		c.JSON(500, gin.H{"status": "Could not connect to DB"})
		return
	}
	defer db.Close()

	// Check if the BookID exists in the DB by querying for the ID
	// Result is scanned into the variable, checkResult
	queryToCheckExistingBook := `SELECT ID FROM BOOKMANAGEMENT WHERE ID=$1;`
	resultToCheckExistingBook := db.QueryRow(queryToCheckExistingBook, addNoteParameters.BookID)
	var checkResult string
	resultToCheckExistingBook.Scan(&checkResult)

	// If the length of checkResult is greater than 0, means the query returned a result, so there is a book by that ID
	// Else, its rejected with a 404 as there is no book by that ID
	if len(checkResult) > 0 {

		//Adds a note to the Book using its ID
		queryToAddANote := `UPDATE BOOKMANAGEMENT SET NOTES = $1 WHERE ID = $2;`
		db.QueryRow(queryToAddANote, sanitizeString(addNoteParameters.Note), addNoteParameters.BookID)
		c.JSON(200, gin.H{"status": "Note added."})

	} else {
		c.JSON(404, gin.H{"status": "No Book with ID, " + addNoteParameters.BookID + " exists"})
	}

}

// Appends a note to a Book, will not clear all the old notes
func addToANote(c *gin.Context) {

	// Variables for DB and Error
	var db *sql.DB
	var err error

	// Creating an instance of the struct, addNoteParameters
	var addNoteParameters AddNoteParameters

	// Bind to the struct's members. If any member is invalid, binding does not happen and an error will be returned. Then its rejected with 400
	if c.BindJSON(&addNoteParameters) != nil {
		c.JSON(400, gin.H{"status": "Incorrect parameters, please provide all required parameters"})
		return
	}

	// Connect to the DB. If there is any issue connecting to the DB, throw a 500 error and return
	db, err = sql.Open("sqlite", "./BOOKMANAGEMENT.db")
	if err != nil {
		c.JSON(500, gin.H{"status": "Could not connect to DB"})
		return
	}
	defer db.Close()

	// Check if the BookID exists in the DB by querying for the ID
	// Result is scanned into the variable, checkResult
	queryToCheckExistingBook := `SELECT ID FROM BOOKMANAGEMENT WHERE ID=$1;`
	resultToCheckExistingBook := db.QueryRow(queryToCheckExistingBook, addNoteParameters.BookID)
	var checkResult string
	resultToCheckExistingBook.Scan(&checkResult)

	// If the length of checkResult is greater than 0, means the query returned a result, so there is a book by that ID
	// Else, its rejected with a 404 as there is no book by that ID
	if len(checkResult) > 0 {

		//Appends a note to the Book's note using its ID
		queryToAddANote := `UPDATE BOOKMANAGEMENT SET NOTES = NOTES || $1 WHERE ID = $2;`
		db.QueryRow(queryToAddANote, " "+sanitizeString(addNoteParameters.Note), addNoteParameters.BookID)
		c.JSON(200, gin.H{"status": "Note appended."})

	} else {
		c.JSON(404, gin.H{"status": "No Book with ID, " + addNoteParameters.BookID + " exists"})
	}

}
