package controller

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"web/src/model"

	"github.com/gin-gonic/gin"
)

// books function is a HTTP handler for GET /books
// returns ALL the books stored in DB as a slice of Book objects
// in users' browser bootstrap will render the JSON in a data table
func books(c *gin.Context) {
	var books []model.Book
	envs := SetEnvs(c)

	DB := model.DB
	rows, err := DB.Query("SELECT * FROM `govwa`.`shelf`")
	if err != nil {
		// debug
		// log.Println(err.Error())
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	defer rows.Close()

	for rows.Next() {
		var book model.Book
		if err := rows.Scan(
			&book.ID,
			&book.Title,
			&book.Author,
			&book.Genre,
			&book.Height,
			&book.Publisher,
			&book.Cover,
		); err != nil {
			// debug
			// log.Println(err.Error())
			c.AbortWithError(http.StatusInternalServerError, err)
			return
		}
		books = append(books, book)
	}
	if err = rows.Err(); err != nil {
		log.Fatalf(err.Error())
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	c.HTML(http.StatusOK, "views/shelf.html", gin.H{
		"title": "GO - Damn Vulnerable Web Application",
		"books": books,
		"envs":  envs,
	})
}

// books function is a HTTP handler for GET /book?q=id
// returns ALL the books stored in DB that match a query
// if parameter is INT than search for id
// if parameter is STRING than search for any other field with LIKE
func book(c *gin.Context) {
	var books []model.Book
	var criteria, qString string
	baseQuery := "SELECT * FROM `govwa`.`shelf` "
	envs := SetEnvs(c)
	qString, success := c.GetQuery("q")
	if !success {
		c.HTML(
			http.StatusOK,
			"views/error.html",
			gin.H{
				"error":        fmt.Errorf("no value supplied for q"),
				"errorMessage": "Empty results",
			},
		)
		return
	}

	qInt, convErr := strconv.Atoi(qString)
	if convErr != nil {
		criteria = fmt.Sprintf("where title like '%%%s%%' or author like '%%%s%%' or genre like '%%%s%%' or publisher = '%%%s%%'", qString, qString, qString, qString)
	} else {
		criteria = fmt.Sprintf("where id = '%d'", qInt)
	}
	query := baseQuery + criteria
	// debug
	// log.Println(query)
	DB := model.DB
	rows, err := DB.Query(query)
	switch {
	case err == sql.ErrNoRows:
		c.HTML(
			http.StatusInternalServerError,
			"views/error.html",
			gin.H{
				"error":        "Searching error",
				"errorMessage": "No book found",
			},
		)
		return
	case err != nil:
		c.HTML(
			http.StatusInternalServerError,
			"views/error.html",
			gin.H{
				"error":        err.Error(),
				"errorMessage": query,
			},
		)
		return
	}
	defer rows.Close()

	for rows.Next() {
		var book model.Book
		if err := rows.Scan(
			&book.ID,
			&book.Title,
			&book.Author,
			&book.Genre,
			&book.Height,
			&book.Publisher,
			&book.Cover,
		); err != nil {
			// debug
			// log.Println(err.Error())
			c.HTML(
				http.StatusInternalServerError,
				"views/error.html",
				gin.H{
					"error":        err.Error(),
					"errorMessage": query,
				},
			)
			return
		}
		// debug
		// log.Printf("Book #%d: %v\n", book.ID, book)
		books = append(books, book)
	}
	if err = rows.Err(); err != nil {
		log.Println(err.Error())
		c.HTML(
			http.StatusInternalServerError,
			"views/error.html",
			gin.H{
				"error":        err.Error(),
				"errorMessage": query,
			},
		)
		return
	}
	c.HTML(http.StatusOK, "views/shelf.html", gin.H{
		"title": "GO - Damn Vulnerable Web Application",
		"books": books,
		"envs":  envs,
	})
}

// booksDetails function is a HTTP handler for GET /book/:id
// it shows page containing details (R/W) about a given book
func bookDetails(c *gin.Context) {
	// debug
	// log.Println("Book details page")
	var book model.Book
	envs := SetEnvs(c)
	bookid := c.Param("id")
	bookIDInt, convErr := strconv.Atoi(bookid)
	if convErr != nil {
		log.Println(convErr.Error())
		c.HTML(
			http.StatusInternalServerError,
			"views/error.html",
			gin.H{
				"error":        convErr.Error(),
				"errorMessage": fmt.Sprintf("Error during conversion: %s", bookid),
			},
		)
		return
	}
	query := fmt.Sprintf("SELECT * FROM `govwa`.`shelf` WHERE id = %d", bookIDInt)
	// DEBUG
	// fmt.Println(query)
	DB := model.DB
	err := DB.QueryRow(query).Scan(&book.ID, &book.Title, &book.Author, &book.Genre, &book.Height, &book.Publisher, &book.Cover)
	switch {
	case err == sql.ErrNoRows:
		c.HTML(
			http.StatusOK,
			"views/error.html",
			gin.H{
				"error":        "Not found",
				"errorMessage": fmt.Sprintf("Profile %s not found", userid),
			},
		)
		return
	case err != nil:
		c.HTML(
			http.StatusOK,
			"views/error.html",
			gin.H{
				"error":        err.Error(),
				"errorMessage": query,
			},
		)
		return
	default:
		c.HTML(http.StatusOK, "views/book.html", gin.H{
			"title": "GO - Damn Vulnerable Web Application",
			"envs":  envs,
			"book":  book,
		})
	}
}
