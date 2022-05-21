package model

import "database/sql"

// Book struct is to store books information
type Book struct {
	ID        int
	Title     sql.NullString
	Author    sql.NullString
	Genre     sql.NullString
	Height    sql.NullString
	Publisher sql.NullString
	Cover     sql.NullString
}
