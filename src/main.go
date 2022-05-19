package main

import (
	"log"
	"web/src/controller"
	"web/src/model"

	_ "github.com/go-sql-driver/mysql"
)

func init() {
	model.Connect()
}

func main() {
	engine := controller.Engine()
	// debug
	log.Println("Server started")
	engine.Run()
}
