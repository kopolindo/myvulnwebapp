package main

import (
	"log"
	"web/controller"
	"web/model"

	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
)

var (
	store cookie.Store
)

func init() {
	//log.Println("connecting to DB")
	model.Connect()
	store = cookie.NewStore([]byte("se45rfgy7yuhji9okmopokmnuygvfr5es2q"))
}

func main() {
	r := gin.New()
	r.Use(gin.Logger())
	r.Use(gin.Recovery())
	r.Use(sessions.Sessions("govwa_cookie", store))

	//Static Files
	r.Static("/css", "./static/css")
	r.Static("/img", "./static/assets/img")
	r.Static("/vendor", "./static/vendor")
	r.Static("/js", "./static/js")
	r.StaticFile("/favicon.ico", "./assets/favicon.ico")

	r.LoadHTMLGlob("templates/**/*")
	controller.Router(r)

	log.Println("Server started")
	r.Run()
}
