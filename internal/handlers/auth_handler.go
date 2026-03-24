package handlers

import (
	"net/http"

	"gosql_interface/internal/database"
	"gosql_interface/internal/models"
	"gosql_interface/internal/utils"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

// LoginPage renders the login page
func LoginPage(c *gin.Context) {
	session := sessions.Default(c)
	userID := session.Get("user_id")

	// If already logged in, redirect to users page
	if userID != nil {
		c.Redirect(http.StatusFound, "/users")
		return
	}

	c.HTML(http.StatusOK, "login.html", gin.H{
		"title": "Login",
	})
}

// LoginHandler handles login form submission
func LoginHandler(c *gin.Context) {
	var loginData struct {
		Login    string `form:"login" json:"login"`
		Password string `form:"password" json:"password"`
	}

	if err := c.ShouldBind(&loginData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid form data"})
		return
	}

	db := database.GetDB()

	// Get user by login
	user, err := models.GetUserByLogin(db, loginData.Login)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid username or password"})
		return
	}

	// Check password
	if err := utils.CheckPassword(user.Password, loginData.Password); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid username or password"})
		return
	}

	// Check if user has admin role
	if user.Role != "admin" {
		c.JSON(http.StatusForbidden, gin.H{"error": "Access denied. Admin privileges required."})
		return
	}

	// Create session
	session := sessions.Default(c)
	session.Set("user_id", user.ID)
	session.Set("user_login", user.Login)
	session.Set("user_role", user.Role)
	if err := session.Save(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save session"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "redirect": "/users"})
}

// LogoutHandler handles logout
func LogoutHandler(c *gin.Context) {
	session := sessions.Default(c)
	session.Clear()
	session.Save()

	c.Redirect(http.StatusFound, "/login")
}
