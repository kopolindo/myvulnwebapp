package controller

import (
	"net/http"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

// This middleware ensures that a request will be aborted with an error
// if the user is not logged in
func ensureLoggedIn(c *gin.Context) {
	loggedInInterface, _ := c.Get("is_logged_in")
	loggedIn := loggedInInterface.(bool)
	if !loggedIn {
		c.HTML(http.StatusUnauthorized, "error.html", gin.H{
			"message":      "Unauthorized",
			"is_logged_in": loggedIn,
		})
		//c.AbortWithStatus(http.StatusUnauthorized)
		c.Abort()
	}
}

// This middleware ensures that a request will be aborted with an error
// if the user is already logged in
func ensureNotLoggedIn(c *gin.Context) {
	loggedInInterface, _ := c.Get("is_logged_in")
	loggedIn := loggedInInterface.(bool)
	if loggedIn {
		c.HTML(http.StatusUnauthorized, "error.html", gin.H{
			"message":      "Unauthorized",
			"is_logged_in": loggedIn,
		})
		c.Abort()
		//c.AbortWithStatus(http.StatusUnauthorized)
	}
}

// setUserStatus set value of loggedFlag (is_logged_in) to true/false
// is_logged_in is used to display or not fields in HTML pages
func setUserStatus(c *gin.Context) {
	c.Set("is_logged_in", false)
	session := sessions.Default(c)
	//session.Options(cookieOptions)
	user := session.Get(userid)
	if user == nil {
		c.Set("is_logged_in", false)
	} else {
		c.Set("is_logged_in", true)
	}
	c.Next()
}
