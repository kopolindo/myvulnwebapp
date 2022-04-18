package controller

import (
	"log"
	"net/http"
	"web/model"

	"github.com/gin-gonic/gin"
)

func shelf(c *gin.Context) {
	var books []model.Book
	setUserStatus(c)
	loggedInInterface, _ := c.Get("is_logged_in")

	DB := model.DB
	rows, err := DB.Query("SELECT * FROM `govwa`.`shelf`")
	if err != nil {
		log.Println(err.Error())
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
			log.Println(err.Error())
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
