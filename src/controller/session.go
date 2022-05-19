package controller

import (
	"fmt"
	"net/http"
	"time"
	"web/src/model"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

type Envs struct {
	LoggedStatus bool
	Email        string
	Role         string
	FirstName    string
	LastName     string
	UserID       int
}

func SetEnvs(c *gin.Context) (env Envs) {
	setUserStatus(c)
	session := sessions.Default(c)
	loggedin, ok := c.Get("is_logged_in")
	if !ok {
		env.LoggedStatus = false
	} else {
		env.LoggedStatus = loggedin.(bool)
	}
	user := session.Get(userid)
	if user == nil {
		env.UserID = -1
	} else {
		env.UserID = user.(int)
	}
	FirstName := session.Get(firstName)
	if FirstName == nil {
		env.FirstName = ""
	} else {
		env.FirstName = FirstName.(string)
	}
	LastName := session.Get(lastName)
	if LastName == nil {
		env.LastName = ""
	} else {
		env.LastName = LastName.(string)
	}
	email := session.Get(userEmail)
	if email == nil {
		env.Email = ""
	} else {
		env.Email = email.(string)
	}
	role := session.Get(userRole)
	if role == nil {
		env.Role = "user"
	} else {
		env.Role = role.(string)
	}
	return
}

// This middleware ensures that a request will be aborted with an error
// if the user is not logged in
func ensureLoggedIn(c *gin.Context) {
	loggedInInterface, _ := c.Get("is_logged_in")
	loggedIn := loggedInInterface.(bool)
	if !loggedIn {
		c.HTML(http.StatusUnauthorized, "views/403.html", gin.H{
			"message":      "ah-ah-ah you didn't say the magic word!",
			"is_logged_in": loggedIn,
		})
		c.Abort()
	}
}

// This middleware ensures that a request will be aborted with an error
// if the user is already logged in
func ensureNotLoggedIn(c *gin.Context) {
	loggedInInterface, _ := c.Get("is_logged_in")
	loggedIn := loggedInInterface.(bool)
	if loggedIn {
		c.HTML(http.StatusUnauthorized, "views/403.html", gin.H{
			"message":      "ah-ah-ah you didn't say the magic word!",
			"is_logged_in": loggedIn,
		})
		c.Abort()
	}
}

// setUserStatus set value of loggedFlag (is_logged_in) to true/false
// is_logged_in is used to display or not fields in HTML pages
func setUserStatus(c *gin.Context) {
	c.Set("is_logged_in", false)
	session := sessions.Default(c)
	//session.Options(cookieOptions)
	user := session.Get(userid)
	role := session.Get(userRole)
	if user == nil {
		c.Set("is_logged_in", false)
	} else {
		c.Set("is_logged_in", true)
	}
	if role == "admin" {
		c.Set("is_admin", true)
	} else {
		c.Set("is_admin", false)
	}
	c.Next()
}

// UpdateActivities write given status for a user identified by its ID
// returns an error
func UpdateActivities(id int, status bool) (sessionError error) {
	var query string
	if status {
		query = "INSERT INTO `govwa`.`activities` (`id`, `last_login`, `status`) values (?, NOW(), 1) ON DUPLICATE KEY UPDATE `last_login` = NOW(), `status` = 1"
	} else {
		query = "INSERT INTO `govwa`.`activities` (`id`, `last_logout`, `status`) values (?, NOW(), 0) ON DUPLICATE KEY UPDATE `last_logout` = NOW(), status = 0"
	}
	// DEBUG
	// fmt.Println(query)
	_, DBError := model.DB.Exec(query, id)
	if DBError != nil {
		return fmt.Errorf("error during activity update")
	}
	setDashboardStatus(true)
	return nil
}

// CheckActivities checks for a given ID its status (login/logout)
// returns status and an error
func CheckActivities(id int) (status bool, sessionError error) {
	var lastLogin, lastLogout time.Time
	DBError := model.DB.QueryRow("SELECT last_login,last_logout,status FROM `govwa`.`activities` WHERE id = ? LIMIT 1", id).Scan(
		&lastLogin,
		&lastLogout,
		&status,
	)
	if DBError != nil {
		return false, fmt.Errorf("error during activity lookup")
	}
	if (!status && (lastLogin.After(lastLogout))) || (status && (lastLogout.After(lastLogin))) {
		return false, fmt.Errorf("error activity record not sounding")
	}
	return status, nil
}
