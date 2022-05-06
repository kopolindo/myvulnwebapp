package controller

import "github.com/gin-contrib/sessions/cookie"

const (
	userid             = "userid"
	userRole           = "userRole"
	userEmail          = "userEmail"
	sessionID          = "sessionID"
	cookieName         = "_govwa"
	uploadPath         = "static/public"
	profilePicturePath = "profiles"
)

var (
	sessionStore = cookie.NewStore([]byte("se45rfgy7yuhji9okmopokmnuygvfr5es2q"))
)
