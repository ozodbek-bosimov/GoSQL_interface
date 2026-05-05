package handlers

import (
	"net/http"

	"gosql_interface/internal/database"
	"gosql_interface/internal/models"
	"gosql_interface/internal/utils"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

func LoginPage(c *gin.Context) {
	session := sessions.Default(c)
	userID := session.Get("user_id")

	if userID != nil {
		c.Redirect(http.StatusFound, "/interfaces")
		return
	}

	c.HTML(http.StatusOK, "login.html", gin.H{
		"title": "Login",
	})
}

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

	user, err := models.GetUserByLogin(db, loginData.Login)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid username or password"})
		return
	}

	if err := utils.CheckPassword(user.Password, loginData.Password); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid username or password"})
		return
	}

	session := sessions.Default(c)
	session.Set("user_id", user.ID)
	session.Set("user_login", user.Login)
	session.Set("user_role", user.Role)
	if err := session.Save(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save session"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "redirect": "/interfaces"})
}

func LogoutHandler(c *gin.Context) {
	session := sessions.Default(c)
	session.Clear()
	session.Save()

	c.Redirect(http.StatusFound, "/login")
}
