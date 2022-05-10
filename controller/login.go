package controller

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"strings"
	"time"
	"web/model"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

type AdminUserLogout struct {
	UserID int `json:"userid" binding:"required"`
}

// loginGet function is a HTTP handler for GET /login
// it returns to the client an HTML page containing login form
func loginGet(c *gin.Context) {
	envs := SetEnvs(c)
	c.HTML(
		http.StatusOK,
		"views/login.html",
		gin.H{
			"title": "GO - Damn Vulnerable Web Application",
			"envs":  envs,
		},
	)
}

// loginPost function is a HTTP handler for POST /login
// it contains authentication logic
func loginPost(c *gin.Context) {
	var emailScan, passwordScan, roleScan, fnameScan, lnameScan string
	var idScan int
	session := sessions.Default(c)
	// Read POST data
	emailIn := strings.Replace(c.PostForm("email"), "'", "\\'", -1)
	passwordIn := strings.Replace(c.PostForm("password"), "'", "\\'", -1)
	// debug
	// log.Println(emailIn, passwordIn)
	qu := fmt.Sprintf("SELECT id FROM govwa.users where email = '%s' limit 1", emailIn)
	errU := model.DB.QueryRow(qu).Scan(&idScan)
	switch {
	case errU == sql.ErrNoRows:
		c.HTML(
			http.StatusOK,
			"views/error.html",
			gin.H{
				"error":        "Credential error",
				"errorMessage": "User not found",
			},
		)
		return
	case errU != nil:
		c.HTML(
			http.StatusOK,
			"views/error.html",
			gin.H{
				"error":        errU.Error(),
				"errorMessage": qu,
			},
		)
		return
	}
	// create query
	qp := fmt.Sprintf(`
		SELECT
			u.id,
			u.email,
			u.password,
			u.role,
			p.first_name,
			p.last_name
		FROM govwa.users as u
		INNER JOIN govwa.profiles as p ON u.id = p.id
		WHERE u.id = '%d' and (u.password = '%s') limit 1`, idScan, passwordIn)
	// debug
	// fmt.Println(qp)
	// qyuery DB
	errP := model.DB.QueryRow(qp).Scan(&idScan, &emailScan, &passwordScan, &roleScan, &fnameScan, &lnameScan)
	switch {
	case errP == sql.ErrNoRows:
		c.HTML(
			http.StatusOK,
			"views/error.html",
			gin.H{
				"error":        "Credential error",
				"errorMessage": "Wrong password",
			},
		)
		return
	case errP != nil:
		c.HTML(
			http.StatusOK,
			"views/error.html",
			gin.H{
				"error":        errP.Error(),
				"errorMessage": qp,
			},
		)
		return
	default:
		c.Set("is_logged_in", true)
		session.Set(userid, idScan)
		session.Set(firstName, fnameScan)
		session.Set(lastName, lnameScan)
		session.Set(userEmail, emailScan)
		session.Set(userRole, roleScan)
		UpdateActivities(idScan, true)
		if err := session.Save(); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save session"})
			return
		}
		location := url.URL{Path: "/api/welcome"}
		c.Redirect(http.StatusFound, location.RequestURI())
	}
}

// logoutGet function is a HTTP handler for GET /logout
// it logs out users and redirects them to /login page
func logoutGet(c *gin.Context) {
	session := sessions.Default(c)

	user := session.Get(userid)
	if user == nil {
		log.Printf("ERROR invalid session token %v\n", user)
		c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("Invalid session token [%v]\n", user)})
		return
	}
	userID := user.(int)
	UpdateActivities(userID, false)
	session.Clear()
	session.Options(sessions.Options{Path: "/", MaxAge: -1})
	if err := session.Save(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save session"})
		return
	}
	location := url.URL{Path: "/login"}
	c.Redirect(http.StatusFound, location.RequestURI())
}

// logoutByAdmin function is a HTTP handler for POST /logout
// it logs out given users and is invoked by admin (on dashboard page)
func logoutByAdmin(c *gin.Context) {
	var user AdminUserLogout
	var lastLogout time.Time
	var firstName string
	c.BindJSON(&user)
	err := UpdateActivities(user.UserID, false)
	if err != nil {
		c.JSON(200, gin.H{"status": err.Error()})
	} else {
		query := fmt.Sprintf(`
			SELECT
				p.first_name,
				a.last_logout
			FROM govwa.activities as a
			INNER JOIN govwa.profiles as p ON a.id = p.id
			WHERE a.id = %d LIMIT 1`,
			user.UserID)
		DBError := model.DB.QueryRow(query).Scan(&firstName, &lastLogout)
		if DBError != nil {
			fmt.Printf("error during activity lookup\n\t%s", DBError.Error())
		}
		c.JSON(200, gin.H{"status": "ok", "firstName": firstName, "lastLogout": lastLogout.String()})
	}
}
