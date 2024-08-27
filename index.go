package main

import (
	"database/sql"

	_ "modernc.org/sqlite"

	"github.com/gin-gonic/gin"
)

func main() {

	request := gin.Default()
	request.GET("/", landingPage)
	request.POST("/addABook", addABook)
	request.GET("/getBookID", getBookID)
	request.Run(":8083")

}

// Landing page route
func landingPage(c *gin.Context) {
	c.JSON(200, "Hello, Welcome to Book Management API")
}

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
		queryToAddABook := `INSERT INTO BOOKMANAGEMENT (ID, BOOK, AUTHOR, TOTALPAGES, READPAGES) Values ($1, $2, $3, $4, $5);`
		db.QueryRow(queryToAddABook, uniqueIDGenerator(), sanitizeString(addABookParameters.BookName), sanitizeString(addABookParameters.AuthorName), addABookParameters.TotalPages, 0)
		c.JSON(200, gin.H{"status": "You have added this Book, " + addABookParameters.BookName})
	}

}

// Defining JSON body for addABook(). It requires 3 JSON key's book, author, totalPages.
type GetBookIDParameters struct {
	BookName   string `form:"book" binding:"required"`
	AuthorName string `form:"author" binding:"required"`
}

// Returns the Book ID
func getBookID(c *gin.Context) {

	// Variables for DB and Error
	var db *sql.DB
	var err error

	// Creating an instance of the struct, GetBookIDParameters
	var getBookIDParameters GetBookIDParameters

	// Bind to the struct's members. If any member is invalid, binding does not happen and an error will be returned. Then its rejected with 400
	if c.Bind(&getBookIDParameters) != nil {
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
	result := db.QueryRow(queryToCheckExistingBook, sanitizeString(getBookIDParameters.BookName), sanitizeString(getBookIDParameters.AuthorName))
	var checkResult string
	result.Scan(&checkResult)

	// If the length of checkResult is greater than 0, means the query returned a result, so there is a book by that author
	// We return the bookname, book author and ID
	// Else, its rejected with a 404 as there is no book by that name and author
	if len(checkResult) > 0 {
		c.JSON(200, gin.H{"bookID": checkResult, "book": sanitizeString(getBookIDParameters.BookName), "author": sanitizeString(getBookIDParameters.AuthorName)})
	} else {
		c.JSON(404, gin.H{"status": "No Book by the name, " + sanitizeString(getBookIDParameters.BookName) + " written by " + sanitizeString(getBookIDParameters.AuthorName) + " exists"})
	}

}
