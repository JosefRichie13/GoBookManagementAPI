# Backend for a Book Management App using Gin Gonic and Go

This repo has the code for a Book Management App Backend. 

The below REST API endpoints are exposed

* GET /getBookID -- Returns a Book's unique ID
  
* GET /getBookDetails -- Returns a Book's details
  
* GET /getAllBooks -- Returns all the available books
  
* GET /getAllUnreadBooks -- Returns all the unread books
  
* GET /getAllReadingBooks -- Returns all the current books being read
  
* GET /getAllFinishedBooks -- Returns all the finished books
  
* GET /getBooksByAuthor -- Returns all the books by an Author
  
* GET /getBooksReadInAPeriod -- Returns all the finished books read in a date range
  
* GET /getBookContaining -- Returns all the books containing a specific word in their name
  
* POST /addABook -- Adds a book
  
* POST /updateBookDetails -- Updates a book's details
  
* POST /startABook -- Starts a book
  
* POST /finishABook -- Finishes a book
  
* POST /updateABook -- Updates a book's read pages
  
* POST /restartABook -- Restarts a finished book
  
* POST /addNote -- Adds a note to a book
  
* POST /addToANote -- Adds to a note of a book
  
* DELETE /deleteBook -- Deletes a book

The entire suite of endpoints with payloads are available in this HAR, [GoBookManagementAPI.har](GoBookManagementAPI.har)
