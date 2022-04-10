package controller

import (
	"net/http"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

// Welcome page, redirect here after successfull login
func me(c *gin.Context) {
	session := sessions.Default(c)
	//session.Options(cookieOptions)
	email := session.Get(userEmail)
	loggedInInterface, _ := c.Get("is_logged_in")
	c.HTML(http.StatusOK, "views/me.html", gin.H{
		"email":        email,
		"is_logged_in": loggedInInterface.(bool),
	})
}
