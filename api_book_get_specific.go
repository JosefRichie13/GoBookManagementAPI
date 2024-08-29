package main

import (
	"database/sql"
	"strconv"

	_ "modernc.org/sqlite"

	"github.com/gin-gonic/gin"
)

// Defining JSON body for getBookID(). It requires 2 Query Parameters book, author.
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

// Defining JSON body for getBookDetails(). It requires 1 Query Parameter bookID.
type GetBookDetailsParameters struct {
	BookID string `form:"bookID" binding:"required"`
}

// Returns a single Book Details
func getBookDetails(c *gin.Context) {

	// Variables for DB and Error
	var db *sql.DB
	var err error

	// Creating an instance of the struct, GetBookDetailsParameters
	var getBookDetailsParameters GetBookDetailsParameters

	// Bind to the struct's members. If any member is invalid, binding does not happen and an error will be returned. Then its rejected with 400
	if c.Bind(&getBookDetailsParameters) != nil {
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

	// Check if the exists in the DB by querying using the ID
	// Result is scanned into the variable, checkResult
	queryToCheckExistingBook := `SELECT * FROM BOOKMANAGEMENT WHERE ID = $1;`
	result := db.QueryRow(queryToCheckExistingBook, getBookDetailsParameters.BookID)

	// Defining a struct to hold all the values from the Query result
	type GetBookDetails struct {
		ID           string
		Book         string
		Author       string
		TotalPages   int
		ReadPages    int
		DateStarted  int
		DateFinished int
		Notes        string
	}

	// Creating an instance of the struct, GetBookDetails
	var getBookDetails GetBookDetails

	// Scan the query result into the struct's members
	result.Scan(&getBookDetails.ID, &getBookDetails.Book, &getBookDetails.Author, &getBookDetails.TotalPages, &getBookDetails.ReadPages,
		&getBookDetails.DateStarted, &getBookDetails.DateFinished, &getBookDetails.Notes)

	// If the length of getBookDetails.ID is greater than 0, means the query returned a result, so there is a book by that ID
	// We return all the details
	// Else, its rejected with a 404 as there is no book by that ID
	if len(getBookDetails.ID) > 0 {
		c.JSON(200, gin.H{"bookID": getBookDetails.ID, "book": getBookDetails.Book, "author": getBookDetails.Author, "totalPages": getBookDetails.TotalPages,
			"readPages": getBookDetails.ReadPages, "dateStarted": convertEpochToDate(getBookDetails.DateStarted), "dateFinished": convertEpochToDate(getBookDetails.DateFinished),
			"notes": getBookDetails.Notes})
	} else {
		c.JSON(404, gin.H{"status": "No Book by ID, " + getBookDetailsParameters.BookID + " exists."})
	}
}

// Defining JSON body for getBooksByAuthor(). It requires 2 Query Parameters book, author.
type GetBooksByAuthorParameters struct {
	AuthorName string `form:"author" binding:"required"`
}

// Returns all Books by a specific author
func getBooksByAuthor(c *gin.Context) {

	// Variables for DB and Error
	var db *sql.DB
	var err error

	// Creating an instance of the struct, GetBookIDParameters
	var getBooksByAuthorParameters GetBooksByAuthorParameters

	// Bind to the struct's members. If any member is invalid, binding does not happen and an error will be returned. Then its rejected with 400
	if c.Bind(&getBooksByAuthorParameters) != nil {
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

	// Query the DB and result is held into the variable, result
	queryToGetAllBooks := `SELECT ID, BOOK FROM BOOKMANAGEMENT WHERE AUTHOR = $1;`
	result, error := db.Query(queryToGetAllBooks, getBooksByAuthorParameters.AuthorName)
	// If there's any error when querying, return it
	if error != nil {
		c.JSON(500, gin.H{"status": "Could not execute Query"})
		return
	}
	defer result.Close()

	// Defining a struct to hold all the values from the Query result
	type GetBookDetails struct {
		ID   string `json:"id"`
		Book string `json:"book"`
	}

	// Creating a slice from the struct
	getBookDetails := []GetBookDetails{}

	// Iterating over the results
	for result.Next() {

		//Creating a new struct
		GetBookDetails := GetBookDetails{}

		// Scan the results into the struct
		result.Scan(&GetBookDetails.ID, &GetBookDetails.Book)

		// Append to the slice
		getBookDetails = append(getBookDetails, GetBookDetails)
	}

	// If there is no result, means, no book by that author exists. Return a 404
	if len(getBookDetails) == 0 {
		c.JSON(404, gin.H{"status": "No book by, " + getBooksByAuthorParameters.AuthorName + " found"})
		return
	}

	// Returning all the data
	c.JSON(200, gin.H{"booksByAuthor": getBookDetails})

}

// Defining JSON body for getBooksReadInAPeriod(). It requires 2 Query Parameters fromDate, toDate.
type GetBooksReadInAPeriodParameters struct {
	FromDate string `form:"fromDate" binding:"required"`
	ToDate   string `form:"toDate" binding:"required"`
}

// Returns all Books by read in a specific period
func getBooksReadInAPeriod(c *gin.Context) {

	// Variables for DB and Error
	var db *sql.DB
	var err error

	// Creating an instance of the struct, GetBooksReadInAPeriodParameters
	var getBooksReadInAPeriodParameters GetBooksReadInAPeriodParameters

	// Bind to the struct's members. If any member is invalid, binding does not happen and an error will be returned. Then its rejected with 400
	if c.Bind(&getBooksReadInAPeriodParameters) != nil {
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

	// Checks if the supplied date is in DD-MMM-YYYY format
	if !checkDateFormat(getBooksReadInAPeriodParameters.FromDate) || !checkDateFormat(getBooksReadInAPeriodParameters.ToDate) {
		c.JSON(400, gin.H{"status": "Incorrect Date format, Date should be in DD-MMM-YYYY format, e.g., 27-Aug-2024"})
		return
	}

	// Query the DB and result is held into the variable, result
	queryToGetAllBooks := `SELECT ID, BOOK, AUTHOR, DATESTARTED, DATEFINISHED FROM BOOKMANAGEMENT WHERE DATESTARTED BETWEEN $1 AND $2 AND DATEFINISHED BETWEEN $1 AND $2;`
	result, error := db.Query(queryToGetAllBooks, convertDateToEpoch(getBooksReadInAPeriodParameters.FromDate), convertDateToEpoch(getBooksReadInAPeriodParameters.ToDate))
	// If there's any error when querying, return it
	if error != nil {
		c.JSON(500, gin.H{"status": "Could not execute Query"})
		return
	}
	defer result.Close()

	// Defining a struct to hold all the values from the Query result
	type GetBookDetails struct {
		ID           string `json:"id"`
		Book         string `json:"book"`
		Author       string `json:"author"`
		DateStarted  string `json:"dateStarted"`
		DateFinished string `json:"dateFinished"`
	}

	// Creating a slice from the struct
	getBookDetails := []GetBookDetails{}

	// Iterating over the results
	for result.Next() {

		//Creating a new struct
		GetBookDetails := GetBookDetails{}

		// Scan the results into the struct
		result.Scan(&GetBookDetails.ID, &GetBookDetails.Book, &GetBookDetails.Author, &GetBookDetails.DateStarted, &GetBookDetails.DateFinished)

		//Converting Date which is in String to Integer and into DD-MMM-YYYY format
		dateStartConversion, errs := strconv.Atoi(GetBookDetails.DateStarted)
		if errs != nil {
			c.JSON(500, gin.H{"status": "Error Processing Start Date"})
			return
		}

		//Converting Date which is in String to Integer and into DD-MMM-YYYY format
		dateFinishConversion, errs := strconv.Atoi(GetBookDetails.DateFinished)
		if errs != nil {
			c.JSON(500, gin.H{"status": "Error Processing Finish Date"})
			return
		}

		//Adding the converted dates
		GetBookDetails.DateStarted = convertEpochToDate(dateStartConversion)
		GetBookDetails.DateFinished = convertEpochToDate(dateFinishConversion)

		// Append to the slice
		getBookDetails = append(getBookDetails, GetBookDetails)
	}

	// If there is no result, means, no book is started and finished between the supplied dates. Return a 404
	if len(getBookDetails) == 0 {
		c.JSON(404, gin.H{"status": "No book read between " + getBooksReadInAPeriodParameters.FromDate + " and " + getBooksReadInAPeriodParameters.ToDate})
		return
	}

	// Returning all the data
	c.JSON(200, gin.H{"booksByAuthor": getBookDetails})

}
