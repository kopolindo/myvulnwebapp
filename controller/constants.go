package controller

import "github.com/gin-contrib/sessions/cookie"

const (
	userid     = "userid"
	userEmail  = "userEmail"
	cookieName = "_govwa"
)

var (
	sessionStore = cookie.NewStore([]byte("se45rfgy7yuhji9okmopokmnuygvfr5es2q"))
)
