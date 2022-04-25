package controller

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"web/model"

	"github.com/gin-gonic/gin"
)

// books function is a HTTP handler for GET /books
// returns ALL the books stored in DB as a slice of Book objects
// in users' browser bootstrap will render the JSON in a data table
func books(c *gin.Context) {
	var books []model.Book
	setUserStatus(c)
	loggedInInterface, _ := c.Get("is_logged_in")

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
		"title":        "GO - Damn Vulnerable Web Application",
		"books":        books,
		"is_logged_in": loggedInInterface.(bool),
	})
}

// books function is a HTTP handler for GET /book/:id
// returns ALL the books stored in DB that match a query
// if parameter is INT than search for id
// if parameter is STRING than search for any other field with LIKE
func book(c *gin.Context) {
	var books []model.Book
	var criteria string
	setUserStatus(c)
	loggedInInterface, _ := c.Get("is_logged_in")

	var ID interface{}
	var err error
	ID, err = strconv.Atoi(c.Param("id"))
	if err != nil {
		ID = c.Param("id")
	}

	baseQuery := "SELECT * FROM `govwa`.`shelf` "

	switch mytype := ID.(type) {
	case string:
		criteria = fmt.Sprintf("where title like '%%%s%%' or author like '%%%s%%' or genre like '%%%s%%' or publisher = '%%%s%%'", ID, ID, ID, ID)
	case int, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64:
		criteria = fmt.Sprintf("where id = '%d'", ID)
	default:
		// debug
		log.Println(mytype)
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
		books = append(books, book)
	}
	if err = rows.Err(); err != nil {
		log.Fatalf(err.Error())
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
		"title":        "GO - Damn Vulnerable Web Application",
		"books":        books,
		"is_logged_in": loggedInInterface.(bool),
	})
}
