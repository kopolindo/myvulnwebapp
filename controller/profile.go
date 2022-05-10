package controller

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"path"
	"path/filepath"
	"web/model"

	"github.com/gin-contrib/sessions"
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
	// fmt.Println(query)
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

// profileUpdateGet handles GET requests for /api/profile/update route
func profileUpdateGet(c *gin.Context) {
	envs := SetEnvs(c)
	var profile model.Profile
	userid := envs.UserID
	query := fmt.Sprintf(`
				SELECT
					u.id,
					p.first_name,
					p.last_name,
					p.image
				FROM
					govwa.users as u
				INNER JOIN govwa.profiles as p ON u.id = p.id
				WHERE
					u.id = '%d'
				LIMIT 1`, userid)
	// DEBUG
	// fmt.Println(query)
	DB := model.DB
	err := DB.QueryRow(query).Scan(&profile.ID, &profile.FirstName, &profile.LastName, &profile.Image)
	switch {
	case err == sql.ErrNoRows:
		c.HTML(
			http.StatusOK,
			"views/error.html",
			gin.H{
				"error":        "Not found",
				"errorMessage": fmt.Sprintf("Profile %d not found", userid),
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
		c.HTML(http.StatusOK, "views/profileUpdate.html", gin.H{
			"title":   "GO - Damn Vulnerable Web Application",
			"envs":    envs,
			"profile": profile,
		})
	}
}

// profileUpdate handles POST requests for /api/profile/:id/update route
// and update profile with information supplied by form
// profile image goes in /static/public/profiles/
func profileUpdate(c *gin.Context) {
	var query string
	session := sessions.Default(c)
	setUserStatus(c)
	envs := SetEnvs(c)
	userid := c.Param("id")
	newFirstName := c.PostForm("firstName")
	if newFirstName == "" {
		newFirstName = envs.FirstName
	}
	newLastName := c.PostForm("lastName")
	if newLastName == "" {
		newLastName = envs.LastName
	}
	file, err := c.FormFile("file")
	if err == nil {
		fname := filepath.Base(file.Filename)
		filename := path.Join(uploadPath, profilePicturePath, fname)
		if err := c.SaveUploadedFile(file, filename); err != nil {
			c.String(http.StatusBadRequest, "upload file err: %s", err.Error())
			return
		}
		query = fmt.Sprintf(`
		UPDATE
			govwa.profiles
		SET
			first_name = "%s",
			last_name = "%s",
			image = "/public/img/%s"
		WHERE
			id = '%s'`, newFirstName, newLastName, fname, userid)
	} else {
		query = fmt.Sprintf(`
		UPDATE
			govwa.profiles
		SET
			first_name = "%s",
			last_name = "%s"
		WHERE
			id = '%s'`, newFirstName, newLastName, userid)
	}
	// DEBUG
	// fmt.Println(query)
	DB := model.DB
	result, e := DB.Exec(query)
	if e != nil {
		// DEBUG
		log.Printf("Error executing query (Exec): %s\n", err.Error())
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
		log.Printf("Error retrieving rows (RowsAffected): %s\n", err.Error())
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
	if newFirstName != "" {
		session.Set(firstName, newFirstName)
	}
	if newLastName != "" {
		session.Set(lastName, newLastName)
	}
	if err := session.Save(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save session"})
		return
	}
	c.Redirect(http.StatusFound, fmt.Sprintf("/api/profile/%s", userid))
}
