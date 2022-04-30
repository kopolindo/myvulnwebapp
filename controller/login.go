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

// loginGet function is a HTTP handler for GET /login
// it returns to the client an HTML page containing login form
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

// loginPost function is a HTTP handler for POST /login
// it contains authentication logic
func loginPost(c *gin.Context) {
	var emailScan, passwordScan, roleScan string
	var idScan int
	session := sessions.Default(c)
	// Read POST data
	emailIn := strings.Replace(c.PostForm("email"), "'", "\\'", -1)
	passwordIn := strings.Replace(c.PostForm("password"), "'", "\\'", -1)
	// debug
	// log.Println(emailIn, passwordIn)
	qu := fmt.Sprintf("SELECT id FROM govwa.users where email = '%s' limit 1", emailIn)
	errU := model.DB.QueryRow(qu).Scan(&idScan)
	switch {
	case errU == sql.ErrNoRows:
		c.HTML(
			http.StatusOK,
			"views/error.html",
			gin.H{
				"error":        "Credential error",
				"errorMessage": "User not found",
			},
		)
		return
	case errU != nil:
		c.HTML(
			http.StatusOK,
			"views/error.html",
			gin.H{
				"error":        errU.Error(),
				"errorMessage": qu,
			},
		)
		return
	}
	// create query
	qp := fmt.Sprintf("SELECT * FROM govwa.users where id = '%d' and (password = '%s') limit 1", idScan, passwordIn)
	// debug
	// fmt.Println(qp)
	// qyuery DB
	errP := model.DB.QueryRow(qp).Scan(&idScan, &emailScan, &passwordScan, &roleScan)
	switch {
	case errP == sql.ErrNoRows:
		c.HTML(
			http.StatusOK,
			"views/error.html",
			gin.H{
				"error":        "Credential error",
				"errorMessage": "Wrong password",
			},
		)
		return
	case errP != nil:
		c.HTML(
			http.StatusOK,
			"views/error.html",
			gin.H{
				"error":        errP.Error(),
				"errorMessage": qp,
			},
		)
		return
	default:
		fmt.Printf("userid: %d\nuserEmail: %s\nuserRole: %s\n", idScan, emailScan, roleScan)
		c.Set("is_logged_in", true)
		session.Set(userid, idScan)
		session.Set(userEmail, emailScan)
		session.Set(userRole, roleScan)
		UpdateActivities(idScan, true)
		if err := session.Save(); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save session"})
			return
		}
		location := url.URL{Path: "/api/me"}
		c.Redirect(http.StatusFound, location.RequestURI())
	}
}

// logoutGet function is a HTTP handler for GET /logout
// it logs out users and redirects them to /login page
func logoutGet(c *gin.Context) {
	session := sessions.Default(c)

	user := session.Get(userid)
	if user == nil {
		log.Printf("ERROR invalid session token %v\n", user)
		c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("Invalid session token [%v]\n", user)})
		return
	}
	userID := user.(int)
	UpdateActivities(userID, false)
	session.Clear()
	session.Options(sessions.Options{Path: "/", MaxAge: -1})
	if err := session.Save(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save session"})
		return
	}
	location := url.URL{Path: "/login"}
	c.Redirect(http.StatusFound, location.RequestURI())
}
