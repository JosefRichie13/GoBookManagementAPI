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
	request.POST("/updateABook", updateABook)
	request.POST("/restartABook", restartABook)
	request.POST("/addNote", addNote)
	request.POST("/addToANote", addToANote)
	request.GET("/getBookID", getBookID)
	request.GET("/getBookDetails", getBookDetails)
	request.GET("/getAllBooks", getAllBooks)
	request.GET("/getAllUnreadBooks", getAllUnreadBooks)
	request.GET("/getAllReadingBooks", getAllReadingBooks)
	request.GET("/getAllFinishedBooks", getAllFinishedBooks)
	request.GET("/getBooksByAuthor", getBooksByAuthor)
	request.GET("/getBooksReadInAPeriod", getBooksReadInAPeriod)
	request.DELETE("/deleteBook", deleteBook)
	request.Run(":8083")

}

// Landing page route
func landingPage(c *gin.Context) {
	c.JSON(200, "Hello, Welcome to Book Management API")
}

// Defining JSON body for deleteBook(). It requires 1 Query Parameter bookID.
type DeleteBookDetailsParameters struct {
	BookID string `form:"bookID" binding:"required"`
}

// Returns a single Book Details
func deleteBook(c *gin.Context) {

	// Variables for DB and Error
	var db *sql.DB
	var err error

	// Creating an instance of the struct, DeleteBookDetailsParameters
	var deleteBookDetailsParameters DeleteBookDetailsParameters

	// Bind to the struct's members. If any member is invalid, binding does not happen and an error will be returned. Then its rejected with 400
	if c.Bind(&deleteBookDetailsParameters) != nil {
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
	// Result is scanned into the variable, result
	queryToCheckExistingBook := `SELECT ID FROM BOOKMANAGEMENT WHERE ID = $1;`
	result := db.QueryRow(queryToCheckExistingBook, deleteBookDetailsParameters.BookID)

	// Defining a struct to hold all the values from the Query result
	type GetBookDetails struct {
		ID string
	}

	// Creating an instance of the struct, GetBookDetails
	var getBookDetails GetBookDetails

	// Scan the query result into the struct's members
	result.Scan(&getBookDetails.ID)

	// If the length of getBookDetails.ID is greater than 0, means the query returned a result, so there is a book by that ID
	// We delete that book
	// Else, its rejected with a 404 as there is no book by that ID
	if len(getBookDetails.ID) > 0 {
		queryToDeleteExistingBook := `DELETE FROM BOOKMANAGEMENT WHERE ID=$1;`
		db.QueryRow(queryToDeleteExistingBook, deleteBookDetailsParameters.BookID)
		c.JSON(200, gin.H{"status": "Book with ID, " + deleteBookDetailsParameters.BookID + " deleted."})

	} else {
		c.JSON(404, gin.H{"status": "No Book by ID, " + deleteBookDetailsParameters.BookID + " exists."})
	}
}
