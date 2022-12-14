package controller

import "github.com/gin-contrib/sessions/cookie"

const (
	userid             = "userid"
	userRole           = "userRole"
	userEmail          = "userEmail"
	firstName          = "firstName"
	lastName           = "lastName"
	sessionID          = "sessionID"
	cookieName         = "_govwa"
	uploadPath         = "static/public"
	profilePicturePath = "profiles"
	bookCoverPath      = "bookCovers"
)

var (
	sessionStore     = cookie.NewStore([]byte("se45rfgy7yuhji9okmopokmnuygvfr5es2q"))
	dashboardChanged = bool(false) // false means "not changed"
)
