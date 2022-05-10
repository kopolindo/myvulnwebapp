package model

import "database/sql"

// Activity struct is to users' activity information
type Activity struct {
	ID         int
	Image      sql.NullString
	Email      sql.NullString
	LastLogin  string
	LastLogout string
	Status     sql.NullInt32
}
