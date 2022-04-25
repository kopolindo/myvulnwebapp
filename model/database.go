package model

import "database/sql"

var DB *sql.DB

// Connect perform connection to the database
// multiStatements true enables stacked queries
// parseTime true enables go to parse datetime as time.Time
func Connect() {
	var err error
	DB, err = sql.Open("mysql", "govwauser:zrXzArJUPyPbB8W@/govwa?multiStatements=true&parseTime=true")
	if err != nil {
		panic(err.Error())
	}
	//defer DB.Close()

	err = DB.Ping()
	if err != nil {
		panic(err.Error())
	}
}
