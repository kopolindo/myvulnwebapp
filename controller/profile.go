package controller

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"path"
	"path/filepath"
	"web/model"

	"github.com/gin-gonic/gin"
)

// Welcome page, redirect here after successfull login
func welcome(c *gin.Context) {
	envs := SetEnvs(c)
	c.HTML(http.StatusOK, "views/welcome.html", gin.H{
		"title": "GO - Damn Vulnerable Web Application",
		"envs":  envs,
	})
}

// Profile page, view only (GET)
func profile(c *gin.Context) {
	envs := SetEnvs(c)
	var profile model.Profile
	userid := c.Param("id")
	query := fmt.Sprintf(`
				SELECT
					u.id,
					u.role,
					p.first_name,
					p.last_name,
					u.email,
					p.image
				FROM
					govwa.users as u
				INNER JOIN govwa.profiles as p ON u.id = p.id
				WHERE
					u.id = '%s'
				LIMIT 1`, userid)
	// DEBUG
	fmt.Println(query)
	DB := model.DB
	err := DB.QueryRow(query).Scan(&profile.ID, &profile.Role, &profile.FirstName, &profile.LastName, &profile.Email, &profile.Image)
	switch {
	case err == sql.ErrNoRows:
		c.HTML(
			http.StatusOK,
			"views/error.html",
			gin.H{
				"error":        "Not found",
				"errorMessage": fmt.Sprintf("Profile %s not found", userid),
			},
		)
		return
	case err != nil:
		c.HTML(
			http.StatusOK,
			"views/error.html",
			gin.H{
				"error":        err.Error(),
				"errorMessage": query,
			},
		)
		return
	default:
		c.HTML(http.StatusOK, "views/profile.html", gin.H{
			"title":   "GO - Damn Vulnerable Web Application",
			"envs":    envs,
			"profile": profile,
		})
	}
}

// imageUpload handles POST requests for /api/profile/:id/upload route
// and upload image files to /static/public/profiles folder
func imageUpload(c *gin.Context) {
	setUserStatus(c)
	//loggedInInterface, _ := c.Get("is_logged_in")
	//session := sessions.Default(c)
	//role := session.Get(userRole)
	userid := c.Param("id")
	file, err := c.FormFile("file")
	if err != nil {
		c.String(http.StatusBadRequest, "get form err: %s", err.Error())
		return
	}
	fname := filepath.Base(file.Filename)
	filename := path.Join(uploadPath, profilePicturePath, fname)
	if err := c.SaveUploadedFile(file, filename); err != nil {
		c.String(http.StatusBadRequest, "upload file err: %s", err.Error())
		return
	}

	query := fmt.Sprintf(`
		UPDATE govwa.profiles
		SET
			image = "/public/img/%s"
		WHERE
			id = '%s'`, fname, userid)
	fmt.Println(query)
	DB := model.DB
	result, e := DB.Exec(query)
	if e != nil {
		log.Println(err.Error())
		c.HTML(
			http.StatusOK,
			"views/error.html",
			gin.H{
				"error":        err.Error(),
				"errorMessage": query,
			},
		)
		return
	}
	rows, err := result.RowsAffected()
	if err != nil {
		c.HTML(
			http.StatusOK,
			"views/error.html",
			gin.H{
				"error":        err.Error(),
				"errorMessage": query,
			},
		)
		return
	}
	if rows != 1 {
		log.Printf("expected to affect 1 row, affected %d\n", rows)
		c.HTML(
			http.StatusOK,
			"views/error.html",
			gin.H{
				"error":        err.Error(),
				"errorMessage": fmt.Sprintf("expected to affect 1 row, affected %d\n", rows),
			},
		)
		return
	}
	c.Redirect(http.StatusFound, fmt.Sprintf("/api/profile/%s", userid))
	/*c.HTML(http.StatusOK, "views/profile.html", gin.H{
		"title":        "GO - Damn Vulnerable Web Application",
		"role":         role,
		"profile":      profile,
		"is_logged_in": loggedInInterface.(bool),
	})*/
}
