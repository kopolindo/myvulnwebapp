package main

import (
	"log"
	"web/controller"
	"web/model"

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
