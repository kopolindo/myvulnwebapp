package controller

import (
	"fmt"
	"log"
	"net/http"
	"time"
	"web/model"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

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

func UpdateActivities(id int, status bool) (sessionError error) {
	var query string
	log.Printf("id (%d) status = %v\t\t", id, status)
	if status {
		query = "INSERT INTO `govwa`.`activities` (`id`, `last_login`, `last_logout`, `status`) values (?, NOW(), NOW(), 1) ON DUPLICATE KEY UPDATE `last_login` = NOW(), `status` = 1"
	} else {
		query = "INSERT INTO `govwa`.`activities` (`id`, `last_logout`, `status`) values (?, NOW(), 0) ON DUPLICATE KEY UPDATE `last_logout` = NOW(), status = 0"
	}
	log.Println(query)
	_, DBError := model.DB.Exec(query, id)
	if DBError != nil {
		return fmt.Errorf("error during activity update")
	}
	return nil
}

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
	return status, nil
}
