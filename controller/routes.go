package controller

import (
	"fmt"
	"net/http"
	"net/url"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

func Engine() *gin.Engine {
	router := gin.New()
	router.Use(gin.Logger())
	router.Use(gin.Recovery())
	// Recovery middleware recovers from any panics and writes a 500 if there was one.
	router.Use(gin.CustomRecovery(func(c *gin.Context, recovered interface{}) {
		if err, ok := recovered.(string); ok {
			c.String(http.StatusInternalServerError, fmt.Sprintf("error: %s", err))
		}
		c.HTML(
			http.StatusInternalServerError,
			"views/error.html",
			gin.H{
				"error":        "Internal Server Error",
				"errorMessage": "Whooopsie daisy! Something went wrong! :_[ ",
			},
		)
	}))
	//Static Files
	router.Static("/css", "./static/css")
	router.Static("/img", "./static/assets/img")
	router.Static("/vendor", "./static/vendor")
	router.Static("/js", "./static/js")
	router.StaticFile("/404", "./templates/views/404.html")
	router.StaticFile("/favicon.ico", "./static/assets/favicon.ico")

	router.LoadHTMLGlob("templates/**/*")
	router.Use(sessions.Sessions(
		cookieName,
		sessionStore,
	))
	// Default routes
	// 404
	router.NoRoute(func(c *gin.Context) {
		location := url.URL{Path: "/404"}
		c.Redirect(http.StatusFound, location.RequestURI())
	})
	router.Use(setUserStatus)
	router.GET("/", index)
	router.GET("/login", ensureNotLoggedIn, loginGet)
	router.POST("/login", ensureNotLoggedIn, loginPost)
	router.GET("/logout", ensureLoggedIn, logoutGet)
	// Private Routes
	private := router.Group("/p")
	private.GET("/me", ensureLoggedIn, me)
	private.GET("/shelf", ensureLoggedIn, shelf)
	return router
}
