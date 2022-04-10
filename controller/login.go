package controller

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"net/url"
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
	var idScan, emailScan, passwordScan string
	session := sessions.Default(c)
	// Read POST data
	emailIn := c.PostForm("email")
	passwordIn := c.PostForm("password")
	// debug
	log.Println(emailIn, passwordIn)
	// create query
	q := fmt.Sprintf("SELECT * FROM govwa.users where email = '%s' and password = '%s' limit 1", emailIn, passwordIn)
	// debug
	fmt.Println(q)
	// qyuery DB
	err := model.DB.QueryRow(q).Scan(&idScan, &emailScan, &passwordScan)
	switch {
	case err == sql.ErrNoRows:
		c.JSON(http.StatusOK, gin.H{"error": "CredentialError"})
		return
	case err != nil:
		c.JSON(http.StatusOK, gin.H{"error": err.Error()})
		return
	default:
		c.Set("is_logged_in", true)
		log.Println(idScan, emailScan, passwordScan)
		session.Set(userid, idScan)
		session.Set(userEmail, emailScan)
		if err := session.Save(); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save session"})
			return
		}
		location := url.URL{Path: "/p/me"}
		c.Redirect(http.StatusFound, location.RequestURI())
	}
}

func logoutHandler(c *gin.Context) {
	session := sessions.Default(c)
	user := session.Get(userid)
	if user == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid session token"})
		return
	}
	session.Delete(userid)
	if err := session.Save(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save session"})
		return
	}
	location := url.URL{Path: "/"}
	c.Redirect(http.StatusFound, location.RequestURI())
	//c.JSON(http.StatusOK, gin.H{"message": "Successfully logged out"})
}
