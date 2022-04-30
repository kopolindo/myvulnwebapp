package controller

import (
	"log"
	"net/http"
	"web/model"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

func dashboard(c *gin.Context) {
	var activities []model.Activity
	setUserStatus(c)
	loggedInInterface, _ := c.Get("is_logged_in")
	session := sessions.Default(c)
	role := session.Get(userRole)

	DB := model.DB
	rows, err := DB.Query("SELECT a.*,u.email FROM `govwa`.`activities` as a INNER JOIN `govwa`.`users` as u ON u.id = a.id")
	if err != nil {
		// debug
		// log.Println(err.Error())
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	defer rows.Close()

	for rows.Next() {
		var activity model.Activity
		if err := rows.Scan(
			&activity.ID,
			&activity.LastLogin,
			&activity.LastLogout,
			&activity.Status,
			&activity.Email,
		); err != nil {
			// debug
			// log.Println(err.Error())
			c.AbortWithError(http.StatusInternalServerError, err)
			return
		}
		activities = append(activities, activity)
	}
	if err = rows.Err(); err != nil {
		log.Fatalf(err.Error())
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	c.HTML(http.StatusOK, "views/dashboard.html", gin.H{
		"title":        "GO - Damn Vulnerable Web Application",
		"books":        activities,
		"role":         role,
		"is_logged_in": loggedInInterface.(bool),
	})
}
