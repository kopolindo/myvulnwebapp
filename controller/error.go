package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func notFound(c *gin.Context) {
	c.HTML(http.StatusOK, "views/error.html", gin.H{
		"message": "Not found",
	})
}
