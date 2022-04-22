package controller

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"strings"
	"web/model"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

func loginGet(c *gin.Context) {
	setUserStatus(c)
	loggedInInterface, _ := c.Get("is_logged_in")
	c.HTML(
		http.StatusOK,
		"views/login.html",
		gin.H{
			"title":        "GO - Damn Vulnerable Web Application",
			"is_logged_in": loggedInInterface.(bool),
		},
	)
}

func loginPost(c *gin.Context) {
	var idScan, emailScan, passwordScan string
	session := sessions.Default(c)
	// Read POST data
	emailIn := strings.Replace(c.PostForm("email"), "'", "\\'", -1)
	passwordIn := strings.Replace(c.PostForm("password"), "'", "\\'", -1)
	// debug
	log.Println(emailIn, passwordIn)
	qu := fmt.Sprintf("SELECT id FROM govwa.users where email = '%s' limit 1", emailIn)
	errU := model.DB.QueryRow(qu).Scan(&idScan)
	switch {
	case errU == sql.ErrNoRows:
		c.HTML(
			http.StatusInternalServerError,
			"views/error.html",
			gin.H{
				"error":        "Credential error",
				"errorMessage": "User not found",
			},
		)
		return
	case errU != nil:
		c.HTML(
			http.StatusInternalServerError,
			"views/error.html",
			gin.H{
				"error":        errU.Error(),
				"errorMessage": qu,
			},
		)
		return
	}
	// create query
	qp := fmt.Sprintf("SELECT * FROM govwa.users where id = '%s' and password = '%s' limit 1", idScan, passwordIn)
	// debug
	fmt.Println(qp)
	// qyuery DB
	errP := model.DB.QueryRow(qp).Scan(&idScan, &emailScan, &passwordScan)
	switch {
	case errP == sql.ErrNoRows:
		c.HTML(
			http.StatusInternalServerError,
			"views/error.html",
			gin.H{
				"error":        "Credential error",
				"errorMessage": "Wrong password",
			},
		)
		return
	case errP != nil:
		c.HTML(
			http.StatusInternalServerError,
			"views/error.html",
			gin.H{
				"error":        errP.Error(),
				"errorMessage": qp,
			},
		)
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

func logoutGet(c *gin.Context) {
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
}
