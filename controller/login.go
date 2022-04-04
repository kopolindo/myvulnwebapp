package controller

import (
	"fmt"
	"log"
	"net/http"
	"web/model"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

func loginGet(c *gin.Context) {
	c.HTML(
		http.StatusOK,
		"views/login.html",
		gin.H{
			"title": "GO - Damn Vulnerable Web Application",
		},
	)
}

func loginPost(c *gin.Context) {
	session := sessions.Default(c)
	user := c.PostForm("email")
	password := c.PostForm("password")
	log.Println(user, password)
	q := fmt.Sprintf("SELECT * FROM govwa.users where email = '%s' and password = '%s'", user, password)
	fmt.Println(q)
	rows, err := model.DB.Query(q)
	if err != nil {
		log.Println(err.Error())
	}

	//check if correct
	for rows.Next() {
		var id, email, password string
		err := rows.Scan(&id, &email, &password)
		if err != nil {
			log.Println(err)
		}
		log.Println(id, email, password)
	}
	err = rows.Err()
	if err != nil {
		log.Println(err)
	}
	//
	c.Set("is_logged_in", true)
	session.Set(userid, user.ID)
	session.Set(userEmail, user.Email)
	if err := session.Save(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save session"})
		return
	}
}
