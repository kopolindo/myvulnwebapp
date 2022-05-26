package model

import (
	"database/sql"
	"html/template"
)

// Book struct is to store books information
type Book struct {
	ID             int
	UnescapedTitle template.HTML
	Title          sql.NullString
	Author         sql.NullString
	Genre          sql.NullString
	Height         sql.NullString
	Publisher      sql.NullString
	Cover          sql.NullString
	BackCover      sql.NullString
}
