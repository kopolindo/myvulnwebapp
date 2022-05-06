package model

import "database/sql"

// Profile struct is to store profile information
type Profile struct {
	ID        int
	Role      sql.NullString
	FirstName sql.NullString
	LastName  sql.NullString
	Email     sql.NullString
	Image     sql.NullString
}
