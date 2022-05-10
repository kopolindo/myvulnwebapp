package controller

import (
	"database/sql"
	"log"
	"net/http"
	"web/model"

	"github.com/gin-gonic/gin"
)

func dashboard(c *gin.Context) {
	var activities []model.Activity
	envs := SetEnvs(c)

	DB := model.DB
	rows, err := DB.Query(`
		SELECT
			a.*,u.email
		FROM
			govwa.activities AS a
		INNER JOIN govwa.users AS u
			ON u.id = a.id`)
	if err != nil {
		// debug
		// log.Println(err.Error())
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	defer rows.Close()
	var lastLoginScan, lastLogoutScan sql.NullTime
	for rows.Next() {
		var activity model.Activity
		if err := rows.Scan(
			&activity.ID,
			lastLoginScan,
			lastLogoutScan,
			&activity.Status,
			&activity.Email,
		); err != nil {
			// debug
			// log.Println(err.Error())
			c.AbortWithError(http.StatusInternalServerError, err)
			return
		}
		activity.LastLogin.Time = lastLoginScan.Time.Format("2006-1-2 15:4:5")
		activity.LastLogout.Time = lastLogoutScan.Time.Format("2006-1-2 15:4:5")
		activities = append(activities, activity)
	}
	if err = rows.Err(); err != nil {
		log.Fatalf(err.Error())
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	c.HTML(http.StatusOK, "views/dashboard.html", gin.H{
		"title": "GO - Damn Vulnerable Web Application",
		"books": activities,
		"envs":  envs,
	})
}
