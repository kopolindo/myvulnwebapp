package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func index(c *gin.Context) {
	setUserStatus(c)
	loggedInInterface, _ := c.Get("is_logged_in")
	c.HTML(
		http.StatusOK,
		"views/index.html",
		gin.H{
			"title":        "GO - Damn Vulnerable Web Application",
			"is_logged_in": loggedInInterface.(bool),
		},
	)
}
