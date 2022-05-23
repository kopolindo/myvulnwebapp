package main

import (
	"log"
	"web/src/controller"
	"web/src/model"
	"web/src/mylog"

	_ "github.com/go-sql-driver/mysql"
)

func init() {
	mylog.Init()
	model.Connect()
}

func main() {
	engine := controller.Engine()
	// debug
	log.Println("Server started")
	engine.Run()
}
