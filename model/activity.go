package model

import "database/sql"

// Activity struct is to users' activity information
type Activity struct {
	ID         int
	Email      sql.NullString
	LastLogin  sql.NullTime
	LastLogout sql.NullTime
	Status     sql.NullInt32
}
