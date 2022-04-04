package controller

import "github.com/gin-gonic/gin"

func Router(r *gin.Engine) {
	r.GET("/", index)
	r.GET("/login", loginGet)
	r.POST("/login", loginPost)
}
