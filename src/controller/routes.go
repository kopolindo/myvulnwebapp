package controller

import (
	"fmt"
	"net/http"
	"net/url"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

// Engine function just set the routes
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
	router.Static("/public/img", "./static/public/profiles")
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
	router.GET("/", func(c *gin.Context) {
		c.Redirect(http.StatusMovedPermanently, "/login")
	})
	router.GET("/login", ensureNotLoggedIn, loginGet)
	router.POST("/login", ensureNotLoggedIn, loginPost)
	router.GET("/logout", ensureLoggedIn, logoutGet)
	router.POST("/logout", ensureLoggedIn, logoutByAdmin)
	// Private Routes
	private := router.Group("/api")
	private.GET("/welcome", ensureLoggedIn, welcome)
	private.GET("/profile/:id", ensureLoggedIn, profile)
	private.GET("/profile/update", ensureLoggedIn, profileUpdateGet)
	private.POST("/profile/:id/update", ensureLoggedIn, profileUpdate)
	private.GET("/books", ensureLoggedIn, books)
	private.GET("/book", ensureLoggedIn, book)
	private.GET("/book/:id", ensureLoggedIn, bookDetails)
	private.GET("/dashboard", ensureLoggedIn, dashboard)
	private.GET("/dashboard/status", ensureLoggedIn, dashboardStatus)
	return router
}
