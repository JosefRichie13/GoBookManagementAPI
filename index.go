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
	request.POST("/updateBookDetails", updateBookDetails)
	request.POST("/startABook", startABook)
	request.POST("/finishABook", finishABook)
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
		generatedID := uniqueIDGenerator()
		queryToAddABook := `INSERT INTO BOOKMANAGEMENT (ID, BOOK, AUTHOR, TOTALPAGES, READPAGES) Values ($1, $2, $3, $4, $5);`
		db.QueryRow(queryToAddABook, generatedID, sanitizeString(addABookParameters.BookName), sanitizeString(addABookParameters.AuthorName), addABookParameters.TotalPages, 0)
		c.JSON(200, gin.H{"status": "Book Added", "bookID": generatedID})
	}

}

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

// Defining JSON body for startABook(). It requires 2 JSON key's bookID, date.
type StartABookParameters struct {
	BookID string `json:"bookID" binding:"required"`
	Date   string `json:"date" binding:"required"`
}

// Starts a Book by updating its DATE STARTED column
func startABook(c *gin.Context) {

	// Variables for DB and Error
	var db *sql.DB
	var err error

	// Creating an instance of the struct, StartABookParameters
	var startABookParameters StartABookParameters

	// Bind to the struct's members. If any member is invalid, binding does not happen and an error will be returned. Then its rejected with 400
	if c.BindJSON(&startABookParameters) != nil {
		c.JSON(400, gin.H{"status": "Incorrect parameters, please provide all required parameters"})
		return
	}

	// Checks if the supplied date is in DD-MMM-YYYY format
	if !checkDateFormat(startABookParameters.Date) {
		c.JSON(400, gin.H{"status": "Incorrect Date format, Date should be in DD-MMM-YYYY format, e.g., 27-Aug-2024"})
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
	resultToCheckExistingBook := db.QueryRow(queryToCheckExistingBook, startABookParameters.BookID)
	var checkResult string
	resultToCheckExistingBook.Scan(&checkResult)

	// If the length of checkResult is greater than 0, means the query returned a result, so there is a book by that ID
	// Else, its rejected with a 404 as there is no book by that ID
	if len(checkResult) > 0 {

		// Check if the BookID exists in the DB by querying for the ID
		// Result is scanned into the variable, result
		queryToCheckExistingBookDates := `SELECT DATESTARTED FROM BOOKMANAGEMENT WHERE ID=$1;`
		result := db.QueryRow(queryToCheckExistingBookDates, startABookParameters.BookID)
		var checkDates int
		result.Scan(&checkDates)

		// If DATESTARTED is not 0, it means that the book is already started
		// We reject wit h a 403
		if checkDates != 0 {
			c.JSON(403, gin.H{"status": "Book with ID, " + startABookParameters.BookID + " is already started"})
			return
		}

		// Update the DATESTARTED, we convert the supplied date in DD-MMM-YYYY format into EpochTime before inserting
		queryToStartABook := `UPDATE BOOKMANAGEMENT SET DATESTARTED = $1 WHERE ID = $2;`
		db.QueryRow(queryToStartABook, convertDateToEpoch(startABookParameters.Date), startABookParameters.BookID)
		c.JSON(200, gin.H{"status": "Book, " + startABookParameters.BookID + " started."})

	} else {
		c.JSON(404, gin.H{"status": "No Book with ID, " + startABookParameters.BookID + " exists"})
	}

}

// Defining JSON body for finishABook(). It requires 2 JSON key's bookID, date.
type FinishABookParameters struct {
	BookID string `json:"bookID" binding:"required"`
	Date   string `json:"date" binding:"required"`
}

// Finishes a Book by updating its DATE FINISHED column
func finishABook(c *gin.Context) {

	// Variables for DB and Error
	var db *sql.DB
	var err error

	// Creating an instance of the struct, FinishABookParameters
	var finishABookParameters FinishABookParameters

	// Bind to the struct's members. If any member is invalid, binding does not happen and an error will be returned. Then its rejected with 400
	if c.BindJSON(&finishABookParameters) != nil {
		c.JSON(400, gin.H{"status": "Incorrect parameters, please provide all required parameters"})
		return
	}

	// Checks if the supplied date is in DD-MMM-YYYY format
	if !checkDateFormat(finishABookParameters.Date) {
		c.JSON(400, gin.H{"status": "Incorrect Date format, Date should be in DD-MMM-YYYY format, e.g., 27-Aug-2024"})
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
	resultToCheckExistingBook := db.QueryRow(queryToCheckExistingBook, finishABookParameters.BookID)
	var checkResult string
	resultToCheckExistingBook.Scan(&checkResult)

	// If the length of checkResult is greater than 0, means the query returned a result, so there is a book by that ID
	// Else, its rejected with a 404 as there is no book by that ID
	if len(checkResult) > 0 {

		// Check if the BookID exists in the DB by querying for the ID
		// Result is scanned into the variable, result
		queryToCheckExistingBookDates := `SELECT TOTALPAGES, DATESTARTED, DATEFINISHED FROM BOOKMANAGEMENT WHERE ID=$1;`
		result := db.QueryRow(queryToCheckExistingBookDates, finishABookParameters.BookID)

		// Defining a struct to hold the data queried by the query and scanning into it
		type CheckDates struct {
			checkDateStarted  int
			checkDateFinished int
			totalPages        int
		}
		var checkDates CheckDates
		result.Scan(&checkDates.totalPages, &checkDates.checkDateStarted, &checkDates.checkDateFinished)

		// If the finished date is less than the started date, reject with 400
		if checkDates.checkDateStarted > convertDateToEpoch(finishABookParameters.Date) {
			c.JSON(400, gin.H{"status": "Finished date cannot be less than Started date"})
			return
		}

		// If DATESTARTED is not 0 and DATEFINISHED is 0, it means that the book is started but not finished
		// We update the DATEFINISHED column to finish the book and set the Read Pages to the Total pages
		// ELSE, the book is not started or is already finished, we reject with 403
		if checkDates.checkDateStarted != 0 && checkDates.checkDateFinished == 0 {
			queryToFinishABook := `UPDATE BOOKMANAGEMENT SET DATEFINISHED = $1, READPAGES =$2 WHERE ID = $3;`
			db.QueryRow(queryToFinishABook, convertDateToEpoch(finishABookParameters.Date), checkDates.totalPages, finishABookParameters.BookID)
			c.JSON(200, gin.H{"status": "Book, " + finishABookParameters.BookID + " finished."})
			return
		} else {
			c.JSON(403, gin.H{"status": "Book, " + finishABookParameters.BookID + " is not started or is finished."})
			return
		}

	} else {
		c.JSON(404, gin.H{"status": "No Book with ID, " + finishABookParameters.BookID + " exists"})
	}

}
