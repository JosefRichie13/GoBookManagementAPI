package main

import (
	"database/sql"

	_ "modernc.org/sqlite"

	"github.com/gin-gonic/gin"
)

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

// Defining JSON body for updateABook(). It requires 2 JSON key's bookID, pages.
type UpdateABookParameters struct {
	BookID string `json:"bookID" binding:"required"`
	Pages  int    `json:"pages" binding:"required"`
}

// Updates a Book by updating its DATE FINISHED column
func updateABook(c *gin.Context) {

	// Variables for DB and Error
	var db *sql.DB
	var err error

	// Creating an instance of the struct, UpdateABookParameters
	var updateABookParameters UpdateABookParameters

	// Bind to the struct's members. If any member is invalid, binding does not happen and an error will be returned. Then its rejected with 400
	if c.BindJSON(&updateABookParameters) != nil {
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
	resultToCheckExistingBook := db.QueryRow(queryToCheckExistingBook, updateABookParameters.BookID)
	var checkResult string
	resultToCheckExistingBook.Scan(&checkResult)

	// If the length of checkResult is greater than 0, means the query returned a result, so there is a book by that ID
	// Else, its rejected with a 404 as there is no book by that ID
	if len(checkResult) > 0 {

		// Check if the BookID exists in the DB by querying for the ID
		// Result is scanned into the variable, result
		queryToCheckExistingBookDates := `SELECT TOTALPAGES, DATESTARTED, DATEFINISHED FROM BOOKMANAGEMENT WHERE ID=$1;`
		result := db.QueryRow(queryToCheckExistingBookDates, updateABookParameters.BookID)

		// Defining a struct to hold the data queried by the query and scanning into it
		type CheckDates struct {
			checkDateStarted  int
			checkDateFinished int
			totalPages        int
		}
		var checkDates CheckDates
		result.Scan(&checkDates.totalPages, &checkDates.checkDateStarted, &checkDates.checkDateFinished)

		// If the suppiled pages is greater less than the total pages, reject with 400
		if updateABookParameters.Pages >= checkDates.totalPages {
			c.JSON(400, gin.H{"status": "Read pages cannot be greater or equal to Total pages."})
			return
		}

		// If DATESTARTED is not 0 and DATEFINISHED is 0, it means that the book is started but not finished
		// We update the Read Pages
		// ELSE, the book is not started or is already finished, we reject with 403
		if checkDates.checkDateStarted != 0 && checkDates.checkDateFinished == 0 {
			queryToFinishABook := `UPDATE BOOKMANAGEMENT SET READPAGES =$1 WHERE ID = $2;`
			db.QueryRow(queryToFinishABook, updateABookParameters.Pages, updateABookParameters.BookID)
			c.JSON(200, gin.H{"status": "Book, " + updateABookParameters.BookID + " updated."})
			return
		} else {
			c.JSON(403, gin.H{"status": "Book, " + updateABookParameters.BookID + " is not started or is finished."})
			return
		}

	} else {
		c.JSON(404, gin.H{"status": "No Book with ID, " + updateABookParameters.BookID + " exists"})
	}

}

// Defining JSON body for restartABook(). It requires 1 JSON key bookID.
type RestartABookParameters struct {
	BookID string `json:"bookID" binding:"required"`
}

// Restarts a Book
func restartABook(c *gin.Context) {

	// Variables for DB and Error
	var db *sql.DB
	var err error

	// Creating an instance of the struct, RestartABookParameters
	var restartABookParameters RestartABookParameters

	// Bind to the struct's members. If any member is invalid, binding does not happen and an error will be returned. Then its rejected with 400
	if c.BindJSON(&restartABookParameters) != nil {
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
	resultToCheckExistingBook := db.QueryRow(queryToCheckExistingBook, restartABookParameters.BookID)
	var checkResult string
	resultToCheckExistingBook.Scan(&checkResult)

	// If the length of checkResult is greater than 0, means the query returned a result, so there is a book by that ID
	// Else, its rejected with a 404 as there is no book by that ID
	if len(checkResult) > 0 {

		// Check if the BookID exists in the DB by querying for the ID
		// Result is scanned into the variable, result
		queryToCheckExistingBookDates := `SELECT DATESTARTED, DATEFINISHED FROM BOOKMANAGEMENT WHERE ID=$1;`
		result := db.QueryRow(queryToCheckExistingBookDates, restartABookParameters.BookID)

		// Defining a struct to hold the data queried by the query and scanning into it
		type CheckDates struct {
			checkDateStarted  int
			checkDateFinished int
		}
		var checkDates CheckDates
		result.Scan(&checkDates.checkDateStarted, &checkDates.checkDateFinished)

		// If DATESTARTED is not 0 and DATEFINISHED is not 0, it means that the book is finished
		// We set READPAGES to 0 and clear out DATESTARTED AND DATEFINISHED, by setting to NULL
		// ELSE, the book is not finished, we reject with 403
		if checkDates.checkDateStarted != 0 && checkDates.checkDateFinished != 0 {
			queryToFinishABook := `UPDATE BOOKMANAGEMENT SET READPAGES = $1, DATESTARTED = 0, DATEFINISHED = 0 WHERE ID = $2;`
			db.QueryRow(queryToFinishABook, 0, restartABookParameters.BookID)
			c.JSON(200, gin.H{"status": "Book, " + restartABookParameters.BookID + " restarted."})
			return
		} else {
			c.JSON(403, gin.H{"status": "Book, " + restartABookParameters.BookID + " is not finished."})
			return
		}

	} else {
		c.JSON(404, gin.H{"status": "No Book with ID, " + restartABookParameters.BookID + " exists"})
	}

}
