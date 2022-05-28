package controller

import (
	"database/sql"
	"fmt"
	"html/template"
	"net/http"
	"path/filepath"
	"strconv"
	"strings"
	"web/src/model"
	"web/src/mylog"

	"github.com/gin-gonic/gin"
)

// books function is a HTTP handler for GET /api/books
// returns ALL the books stored in DB as a slice of Book objects
// in users' browser bootstrap will render the JSON in a data table
func books(c *gin.Context) {
	var books []model.Book
	envs := SetEnvs(c)

	DB := model.DB
	rows, err := DB.Query("SELECT * FROM `govwa`.`shelf`")
	if err != nil {
		// debug
		mylog.Debug.Println(err.Error())
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
			&book.BackCover,
		); err != nil {
			// debug
			mylog.Debug.Println(err.Error())
			c.AbortWithError(http.StatusInternalServerError, err)
			return
		}
		books = append(books, book)
	}
	if err = rows.Err(); err != nil {
		mylog.Error.Println(err.Error())
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	c.HTML(http.StatusOK, "views/shelf.html", gin.H{
		"title": "GO - Damn Vulnerable Web Application",
		"books": books,
		"envs":  envs,
	})
}

// books function is a HTTP handler for GET /api/book?q=id
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
	mylog.Debug.Println(query)
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
			&book.BackCover,
		); err != nil {
			// debug
			mylog.Debug.Println(err.Error())
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
		mylog.Debug.Printf("Book #%d: %v\n", book.ID, book)
		books = append(books, book)
	}
	if err = rows.Err(); err != nil {
		mylog.Error.Println(err.Error())
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

// bookDetailsUpdate handles POST requests for /api/book/:id/update route
// and update book details with information supplied by form
// cover images goes in /static/public/bookCovers/

func bookDetailsUpdate(c *gin.Context) {
	// debug
	mylog.Debug.Println("Update book details")
	var query, values string
	//envs := SetEnvs(c)
	bookid := c.Param("id")
	bookIDInt, convErr := strconv.Atoi(bookid)
	if convErr != nil {
		mylog.Error.Println(convErr.Error())
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
	// Handle form input
	inputTitle := c.PostForm("inputTitle")
	if inputTitle != "" {
		values = fmt.Sprintf("title = '%s'", inputTitle)
	}
	inputAuthor := c.PostForm("inputAuthor")
	if inputAuthor != "" {
		if values != "" {
			values = fmt.Sprintf("%s,author = '%s'", values, inputAuthor)
		} else {
			values = fmt.Sprintf("author = '%s'", inputAuthor)
		}
	}
	inputGenre := c.PostForm("inputGenre")
	if inputGenre != "" {
		if values != "" {
			values = fmt.Sprintf("%s,genre = '%s'", values, inputGenre)
		} else {
			values = fmt.Sprintf("genre = '%s'", inputGenre)
		}
	}
	inputPublisher := c.PostForm("inputPublisher")
	if inputPublisher != "" {
		if values != "" {
			values = fmt.Sprintf("%s,publisher = '%s'", values, inputPublisher)
		} else {
			values = fmt.Sprintf("publisher = '%s'", inputPublisher)
		}
	}
	inputBackCover := c.PostForm("inputBackCover")
	if inputBackCover != "" {
		if values != "" {
			values = fmt.Sprintf("%s,back_cover = '%s'", values, inputBackCover)
		} else {
			values = fmt.Sprintf("back_cover = '%s'", inputBackCover)
		}
	}
	file, err := c.FormFile("inputCover")
	mylog.Debug.Println("== FILENAME ==", file.Filename)
	if err == nil {
		fname := filepath.Base(file.Filename)
		if err := c.SaveUploadedFile(file, file.Filename); err != nil {
			c.String(http.StatusBadRequest, "upload file err: %s", err.Error())
			return
		}
		if values != "" {
			values = fmt.Sprintf("%s,cover = '/public/covers/%s'", values, fname)
		} else {
			values = fmt.Sprintf("cover = '/public/covers/%s'", fname)
		}
	}
	// Query creation
	query = fmt.Sprintf(`
		UPDATE
			govwa.shelf
		SET
			%s
		WHERE
			id = '%d'`, values, bookIDInt)
	// DEBUG
	mylog.Debug.Println(query)
	DB := model.DB
	result, err := DB.Exec(query)
	if err != nil {
		// DEBUG
		mylog.Error.Printf("Error executing query (Exec): %s\n", err.Error())
		c.HTML(
			http.StatusOK,
			"views/error.html",
			gin.H{
				"error":        err.Error(),
				"errorMessage": query,
			},
		)
		return
	}
	rows, err := result.RowsAffected()
	if err != nil {
		mylog.Error.Printf("Error retrieving rows (RowsAffected): %s\n", err.Error())
		c.HTML(
			http.StatusOK,
			"views/error.html",
			gin.H{
				"error":        err.Error(),
				"errorMessage": query,
			},
		)
		return
	}
	if rows != 1 {
		mylog.Error.Printf("== QUERY == %s\n", query)
		mylog.Error.Printf("expected to affect 1 row, affected %d\n", rows)
		c.HTML(
			http.StatusOK,
			"views/error.html",
			gin.H{
				"error":        err.Error(),
				"errorMessage": fmt.Sprintf("expected to affect 1 row, affected %d\n", rows),
			},
		)
		return
	}
	c.Redirect(http.StatusFound, fmt.Sprintf("/api/book/%d", bookIDInt))
}

// booksDetails function is a HTTP handler for GET /book/:id
// it shows page containing details (R/W) about a given book
func bookDetails(c *gin.Context) {
	// debug
	mylog.Debug.Println("Book details page")
	var book model.Book
	envs := SetEnvs(c)
	bookid := c.Param("id")
	bookIDInt, convErr := strconv.Atoi(bookid)
	if convErr != nil {
		mylog.Error.Println(convErr.Error())
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
	mylog.Debug.Println(query)
	DB := model.DB
	err := DB.QueryRow(query).Scan(
		&book.ID,
		&book.Title,
		&book.Author,
		&book.Genre,
		&book.Height,
		&book.Publisher,
		&book.Cover,
		&book.BackCover,
	)
	book.UnescapedTitle = template.HTML(strings.Replace(book.Title.String, " ", "<br>", 1))
	fmt.Printf("TITLE: %s\n", book.UnescapedTitle)
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
