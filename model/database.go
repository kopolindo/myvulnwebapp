package model

import "database/sql"

var DB *sql.DB

func Connect() {
	var err error
	DB, err = sql.Open("mysql", "govwauser:zrXzArJUPyPbB8W@/govwa")
	if err != nil {
		panic(err.Error())
	}
	//defer DB.Close()

	err = DB.Ping()
	if err != nil {
		panic(err.Error())
	}
}
