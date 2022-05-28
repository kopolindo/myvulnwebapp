package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func debug(c *gin.Context) {
	envs := SetEnvs(c)
	c.HTML(http.StatusOK, "views/message.html", gin.H{
		"title":   "GO - Damn Vulnerable Web Application",
		"envs":    envs,
		"message": "IT'S WORKING!!",
	})
}
