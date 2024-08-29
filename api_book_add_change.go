package main

import (
	"database/sql"

	_ "modernc.org/sqlite"

	"github.com/gin-gonic/gin"
)

// Defining JSON body for addABook(). It requires 3 JSON key's book, author, totalPages.
type AddABookParameters struct {
	BookName   string `json:"book" binding:"required"`
	AuthorName string `json:"author" binding:"required"`
	TotalPages int    `json:"totalPages" binding:"required"`
}

// Adds a Book to the DB
func addABook(c *gin.Context) {

	// Variables for DB and Error
	var db *sql.DB
	var err error

	// Creating an instance of the struct, AddABookParameters
	var addABookParameters AddABookParameters

	// Bind to the struct's members. If any member is invalid, binding does not happen and an error will be returned. Then its rejected with 400
	if c.Bind(&addABookParameters) != nil {
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

	// Check if the Book and the Author exists in the DB by querying for the ID
	// Result is scanned into the variable, checkResult
	queryToCheckExistingBook := `SELECT ID FROM BOOKMANAGEMENT WHERE BOOK=$1 AND AUTHOR=$2;`
	result := db.QueryRow(queryToCheckExistingBook, sanitizeString(addABookParameters.BookName), sanitizeString(addABookParameters.AuthorName))
	var checkResult string
	result.Scan(&checkResult)

	// If the length of checkResult is greater than 0, means the query returned a result, so there is a book by that author already, then its rejected with a 403
	// Else, its added to the DB
	if len(checkResult) > 0 {
		c.JSON(403, gin.H{"status": "Book, " + addABookParameters.BookName + " by " + addABookParameters.AuthorName + " already exists"})
	} else {
		generatedID := uniqueIDGenerator()
		queryToAddABook := `INSERT INTO BOOKMANAGEMENT (ID, BOOK, AUTHOR, TOTALPAGES, READPAGES, DATESTARTED, DATEFINISHED, NOTES) Values ($1, $2, $3, $4, $5, $6, $7, $8);`
		db.QueryRow(queryToAddABook, generatedID, sanitizeString(addABookParameters.BookName), sanitizeString(addABookParameters.AuthorName),
			addABookParameters.TotalPages, 0, 0, 0, "")
		c.JSON(200, gin.H{"status": "Book Added", "bookID": generatedID})
	}

}

// Defining JSON body for updateBookDetails(). It requires 4 JSON key's bookID, book, author and totalPages.
type UpdateBookDetailsParameters struct {
	BookID     string `json:"bookID" binding:"required"`
	BookName   string `json:"book" binding:"required"`
	AuthorName string `json:"author" binding:"required"`
	TotalPages int    `json:"totalPages" binding:"required"`
}

// Updates an existing Book's details, Name, Author and Total Pages
func updateBookDetails(c *gin.Context) {

	// Variables for DB and Error
	var db *sql.DB
	var err error

	// Creating an instance of the struct, UpdateBookDetailsParameters
	var updateBookDetailsParameters UpdateBookDetailsParameters

	// Bind to the struct's members. If any member is invalid, binding does not happen and an error will be returned. Then its rejected with 400
	if c.BindJSON(&updateBookDetailsParameters) != nil {
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
	result := db.QueryRow(queryToCheckExistingBook, updateBookDetailsParameters.BookID)
	var checkResult string
	result.Scan(&checkResult)

	// If the length of checkResult is greater than 0, means the query returned a result, so there is a book by that author
	// Else, its rejected with a 404 as there is no book by that ID
	if len(checkResult) > 0 {

		// Check if the update book details match any exisiting book details in the DB
		// If yes, reject with 403
		queryToCheckIfBookExists := `SELECT ID FROM BOOKMANAGEMENT WHERE BOOK=$1 AND AUTHOR=$2 AND TOTALPAGES=$3;`
		resultToCheckIfBookExists := db.QueryRow(queryToCheckIfBookExists, sanitizeString(updateBookDetailsParameters.BookName), sanitizeString(updateBookDetailsParameters.AuthorName), updateBookDetailsParameters.TotalPages)
		var checkIfBookExists string
		resultToCheckIfBookExists.Scan(&checkIfBookExists)
		if len(checkIfBookExists) > 0 {
			c.JSON(403, gin.H{"status": "Same Book by the same author with the same page number already exists."})
			return
		}

		// Then if the update book details are different, update the book details
		queryToUpdateABook := `UPDATE BOOKMANAGEMENT SET BOOK = $1, AUTHOR = $2, TOTALPAGES =$3 WHERE ID = $4;`
		db.QueryRow(queryToUpdateABook, sanitizeString(updateBookDetailsParameters.BookName), sanitizeString(updateBookDetailsParameters.AuthorName), updateBookDetailsParameters.TotalPages, updateBookDetailsParameters.BookID)
		c.JSON(200, gin.H{"status": "Book, " + updateBookDetailsParameters.BookID + " updated."})

	} else {
		c.JSON(404, gin.H{"status": "No Book with ID, " + updateBookDetailsParameters.BookID + " exists"})
	}

}
