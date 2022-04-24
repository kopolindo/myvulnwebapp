package controller

import (
	"net/http"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

// Welcome page, redirect here after successfull login
func me(c *gin.Context) {
	setUserStatus(c)
	loggedInInterface, _ := c.Get("is_logged_in")
	session := sessions.Default(c)
	//session.Options(cookieOptions)
	email := session.Get(userEmail)
	role := session.Get(userRole)
	c.HTML(http.StatusOK, "views/me.html", gin.H{
		"title":        "GO - Damn Vulnerable Web Application",
		"email":        email,
		"role":         role,
		"is_logged_in": loggedInInterface.(bool),
	})
}
