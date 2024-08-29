package main

import (
	"database/sql"
	"fmt"
	"math"
	"strconv"

	_ "modernc.org/sqlite"

	"github.com/gin-gonic/gin"
)

// Returns all the available Book Details
func getAllBooks(c *gin.Context) {

	// Variables for DB and Error
	var db *sql.DB
	var err error

	// Connect to the DB. If there is any issue connecting to the DB, throw a 500 error and return
	db, err = sql.Open("sqlite", "./BOOKMANAGEMENT.db")
	if err != nil {
		c.JSON(500, gin.H{"status": "Could not connect to DB"})
		return
	}
	defer db.Close()

	// Query the DB and result is held into the variable, result
	queryToGetAllBooks := `SELECT * FROM BOOKMANAGEMENT;`
	result, error := db.Query(queryToGetAllBooks)
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
		TotalPages   int    `json:"totalPages"`
		ReadPages    int    `json:"readPages"`
		DateStarted  string `json:"dateStarted"`
		DateFinished string `json:"dateFinished"`
		Notes        string `json:"notes"`
	}

	// Creating a slice from the struct
	getBookDetails := []GetBookDetails{}

	// Iterating over the results
	for result.Next() {

		//Creating a new struct
		GetBookDetails := GetBookDetails{}
		// Scan the results into the struct
		result.Scan(&GetBookDetails.ID, &GetBookDetails.Book, &GetBookDetails.Author, &GetBookDetails.TotalPages, &GetBookDetails.ReadPages,
			&GetBookDetails.DateStarted, &GetBookDetails.DateFinished, &GetBookDetails.Notes)

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

	// Returning all the data
	c.JSON(200, gin.H{"allBookDetails": getBookDetails})

}

// Returns all the Unread Book Details
func getAllUnreadBooks(c *gin.Context) {

	// Variables for DB and Error
	var db *sql.DB
	var err error

	// Connect to the DB. If there is any issue connecting to the DB, throw a 500 error and return
	db, err = sql.Open("sqlite", "./BOOKMANAGEMENT.db")
	if err != nil {
		c.JSON(500, gin.H{"status": "Could not connect to DB"})
		return
	}
	defer db.Close()

	// Query the DB and result is held into the variable, result
	queryToGetAllBooks := `SELECT ID, BOOK, AUTHOR FROM BOOKMANAGEMENT where DATESTARTED IS 0 AND DATEFINISHED IS 0;`
	result, error := db.Query(queryToGetAllBooks)
	// If there's any error when querying, return it
	if error != nil {
		c.JSON(500, gin.H{"status": "Could not execute Query"})
		return
	}
	defer result.Close()

	// Defining a struct to hold all the values from the Query result
	type GetBookDetails struct {
		ID     string `json:"id"`
		Book   string `json:"book"`
		Author string `json:"author"`
	}

	// Creating a slice from the struct
	getBookDetails := []GetBookDetails{}

	// Iterating over the results
	for result.Next() {

		//Creating a new struct
		GetBookDetails := GetBookDetails{}
		// Scan the results into the struct
		result.Scan(&GetBookDetails.ID, &GetBookDetails.Book, &GetBookDetails.Author)
		// Append to the slice
		getBookDetails = append(getBookDetails, GetBookDetails)
	}

	// Returning all the data
	c.JSON(200, gin.H{"unreadBookDetails": getBookDetails})

}

// Returns all current Reading Book Details
func getAllReadingBooks(c *gin.Context) {

	// Variables for DB and Error
	var db *sql.DB
	var err error

	// Connect to the DB. If there is any issue connecting to the DB, throw a 500 error and return
	db, err = sql.Open("sqlite", "./BOOKMANAGEMENT.db")
	if err != nil {
		c.JSON(500, gin.H{"status": "Could not connect to DB"})
		return
	}
	defer db.Close()

	// Query the DB and result is held into the variable, result
	queryToGetAllBooks := `SELECT ID, BOOK, AUTHOR, DATESTARTED, TOTALPAGES, READPAGES FROM BOOKMANAGEMENT where DATESTARTED IS NOT 0 AND DATEFINISHED IS 0;`
	result, error := db.Query(queryToGetAllBooks)
	// If there's any error when querying, return it
	if error != nil {
		c.JSON(500, gin.H{"status": "Could not execute Query"})
		return
	}
	defer result.Close()

	// Defining a struct to hold all the values from the Query result
	type GetBookDetails struct {
		ID                 string  `json:"id"`
		Book               string  `json:"book"`
		Author             string  `json:"author"`
		DateStarted        string  `json:"dateStarted"`
		TotalPages         int     `json:"totalPages"`
		ReadPages          int     `json:"readPages"`
		RemainingPages     int     `json:"remainingPages"`
		PercentageFinished float64 `json:"percentageFinished"`
		PercentageLeft     float64 `json:"percentageLeft"`
	}

	// Creating a slice from the struct
	getBookDetails := []GetBookDetails{}

	// Iterating over the results
	for result.Next() {

		//Creating a new struct
		GetBookDetails := GetBookDetails{}

		// Scan the results into the struct
		result.Scan(&GetBookDetails.ID, &GetBookDetails.Book, &GetBookDetails.Author, &GetBookDetails.DateStarted, &GetBookDetails.TotalPages, &GetBookDetails.ReadPages)

		//Converting Date which is in String to Integer and into DD-MMM-YYYY format
		dateStartConversion, errs := strconv.Atoi(GetBookDetails.DateStarted)
		if errs != nil {
			c.JSON(500, gin.H{"status": "Error Processing Start Date"})
			return
		}

		//Adding the converted dates
		GetBookDetails.DateStarted = convertEpochToDate(dateStartConversion)

		// Calculating remaining pages, subtracting Read pages from Total pages, gives us the remaining pages
		GetBookDetails.RemainingPages = GetBookDetails.TotalPages - GetBookDetails.ReadPages

		// Calculating the percentage of pages left to read
		GetBookDetails.PercentageLeft = 100 - (float64(GetBookDetails.ReadPages)/float64(GetBookDetails.TotalPages))*100
		GetBookDetails.PercentageLeft = math.Round(GetBookDetails.PercentageLeft)

		// Calculating the percentage of finished pages
		GetBookDetails.PercentageFinished = (float64(GetBookDetails.ReadPages) / float64(GetBookDetails.TotalPages)) * 100
		GetBookDetails.PercentageFinished = math.Round(GetBookDetails.PercentageFinished)

		// Append to the slice
		getBookDetails = append(getBookDetails, GetBookDetails)
	}

	// Returning all the data
	c.JSON(200, gin.H{"currentlyReadingBooks": getBookDetails})

}

// Returns all Finished Book Details
func getAllFinishedBooks(c *gin.Context) {

	// Variables for DB and Error
	var db *sql.DB
	var err error

	// Connect to the DB. If there is any issue connecting to the DB, throw a 500 error and return
	db, err = sql.Open("sqlite", "./BOOKMANAGEMENT.db")
	if err != nil {
		c.JSON(500, gin.H{"status": "Could not connect to DB"})
		return
	}
	defer db.Close()

	// Query the DB and result is held into the variable, result
	queryToGetAllBooks := `SELECT ID, BOOK, AUTHOR, DATESTARTED, DATEFINISHED FROM BOOKMANAGEMENT where DATESTARTED IS NOT 0 AND DATEFINISHED IS NOT 0;`
	result, error := db.Query(queryToGetAllBooks)
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
		DaysRead     int64  `json:"daysRead"`
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

		// Calculating the days in which a book was read
		// Convert Start and Finished dates, in Epoch times, to Integers from strings
		startDateInInt, startErr := strconv.ParseInt(GetBookDetails.DateStarted, 10, 64)
		if startErr != nil {
			fmt.Println("Error converting Start Date to int")
		}
		finishDateInInt, finishErr := strconv.ParseInt(GetBookDetails.DateFinished, 10, 64)
		if finishErr != nil {
			fmt.Println("Error converting Finish Date to int")
		}

		// Calculate no of days read by subtracting Start date in Epoch from Finished date in Epoch and dividing it by 86400
		// As Epoch is in seconds, dividing by 86400 (no of seconds in a day) will give the no of days
		// If the daysRead is 0, happens when the Date Start and Date Finished are the same, we set it to 1
		daysRead := (finishDateInInt - startDateInInt) / 86400
		if daysRead == 0 {
			GetBookDetails.DaysRead = 1
		} else {
			GetBookDetails.DaysRead = daysRead
		}

		//Adding the converted dates
		GetBookDetails.DateStarted = convertEpochToDate(dateStartConversion)
		GetBookDetails.DateFinished = convertEpochToDate(dateFinishConversion)

		// Append to the slice
		getBookDetails = append(getBookDetails, GetBookDetails)
	}

	// Returning all the data
	c.JSON(200, gin.H{"finishedBookDetails": getBookDetails})

}
