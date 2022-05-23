package model

import (
	"database/sql"
	"time"
	"web/src/mylog"
)

var DB *sql.DB

// Connect perform connection to the database
// multiStatements true enables stacked queries
// parseTime true enables go to parse datetime as time.Time
func Connect() {
	var err error

	for {
		DB, err = sql.Open("mysql", "govwauser:zrXzArJUPyPbB8W@tcp(mariadb-govwa:3306)/govwa?multiStatements=true&parseTime=true")
		if err != nil {
			mylog.Error.Println(err.Error())
		}
		//defer DB.Close()

		err = DB.Ping()
		if err == nil {
			mylog.Info.Println("CONNECTED! GO ON!")
			break
		}
		time.Sleep(1 * time.Second)
		continue
	}

}
