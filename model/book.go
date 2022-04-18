package model

import "database/sql"

type Book struct {
	ID        int
	Title     sql.NullString
	Author    sql.NullString
	Genre     sql.NullString
	Height    sql.NullString
	Publisher sql.NullString
}
