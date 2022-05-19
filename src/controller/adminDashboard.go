package controller

import (
	"database/sql"
	"log"
	"net/http"
	"web/src/model"

	"github.com/gin-gonic/gin"
)

// dashboard function handles GET to /api/dashboard route
// returns a model.Activity slice containing
//		ID, email, lastLogin, lastLogout, status (logged in yes/no), force logout button
func dashboard(c *gin.Context) {
	envs := SetEnvs(c)
	DB := model.DB
	var lastLoginScan, lastLogoutScan sql.NullTime
	var activities []model.Activity
	var layout = "2006-01-02 15:04:05"
	var query = `SELECT
					a.id,
					a.last_login,
					a.last_logout,
					a.status,
					u.email,
					p.image
				FROM
					govwa.activities AS a
				INNER JOIN govwa.users AS u
					ON u.id = a.id
				INNER JOIN govwa.profiles AS p
					ON p.id = a.id`
	rows, err := DB.Query(query)
	if err != nil {
		// debug
		// log.Println(query)
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	defer rows.Close()
	for rows.Next() {
		var activity model.Activity
		if err := rows.Scan(
			&activity.ID,
			&lastLoginScan,
			&lastLogoutScan,
			&activity.Status,
			&activity.Email,
			&activity.Image,
		); err != nil {
			// debug
			// log.Println(query)
			c.AbortWithError(http.StatusInternalServerError, err)
			return
		}

		activity.LastLogin = lastLoginScan.Time.Format(layout)
		activity.LastLogout = lastLogoutScan.Time.Format(layout)
		activities = append(activities, activity)
	}
	if err = rows.Err(); err != nil {
		log.Fatalf(err.Error())
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	c.HTML(http.StatusOK, "views/dashboard.html", gin.H{
		"title":      "GO - Damn Vulnerable Web Application",
		"activities": activities,
		"envs":       envs,
	})
}

// setDashboardStatus functions set the value of dashboardChanged variable
func setDashboardStatus(value bool) {
	dashboardChanged = value
}

// getDashboardStatus functions returns value of dashboardChanged variable
func getDashboardStatus() bool {
	return dashboardChanged
}

// dashboardStatus functions handles GET /dashboard/status route
// returns dashboardChanged true/false
func dashboardStatus(c *gin.Context) {
	status := getDashboardStatus()
	if status {
		setDashboardStatus(false)
	}
	c.JSON(http.StatusOK, gin.H{"status": status})
}
