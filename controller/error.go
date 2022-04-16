package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func notFound(c *gin.Context) {
	setUserStatus(c)
	loggedInInterface, _ := c.Get("is_logged_in")
	c.HTML(http.StatusOK, "views/error.html", gin.H{
		"is_logged_in": loggedInInterface.(bool),
		"message":      "Not found",
	})
}
