package main

import (
	"log"
	"web/controller"
	"web/model"

	"github.com/gin-contrib/sessions/cookie"
	_ "github.com/go-sql-driver/mysql"
)

var (
	store cookie.Store
)

func init() {
	//log.Println("connecting to DB")
	model.Connect()
}

func main() {
	engine := controller.Engine()

	log.Println("Server started")
	engine.Run()
}
