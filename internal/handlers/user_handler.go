package handlers

import (
	"net/http"
	"strconv"

	"gosql_interface/internal/database"
	"gosql_interface/internal/models"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

func ListUsersPage(c *gin.Context) {
	session := sessions.Default(c)
	c.HTML(http.StatusOK, "users.html", gin.H{
		"title":      "Users",
		"activePage": "users",
		"userLogin":  session.Get("user_login"),
		"userRole":   session.Get("user_role"),
	})
}

func GetUsersAPI(c *gin.Context) {
	search := c.DefaultQuery("search", "")
	db := database.GetDB()

	users, err := models.GetAllUsers(db, search)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	for i := range users {
		users[i].Password = ""
	}

	c.JSON(http.StatusOK, users)
}

func CreateUserAPI(c *gin.Context) {
	var user models.User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid data"})
		return
	}

	db := database.GetDB()
	if err := models.CreateUser(db, &user); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	user.Password = ""
	c.JSON(http.StatusCreated, user)
}

func UpdateUserAPI(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	var user models.User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid data"})
		return
	}

	user.ID = id
	db := database.GetDB()
	if err := models.UpdateUser(db, &user); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	user.Password = ""
	c.JSON(http.StatusOK, user)
}

func DeleteUserAPI(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	db := database.GetDB()
	if err := models.DeleteUser(db, id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "User deleted successfully"})
}

func GetUserSchemaAPI(c *gin.Context) {
	schema := []map[string]interface{}{
		{"name": "id", "type": "BIGINT", "nullable": false, "primary": true},
		{"name": "name", "type": "TEXT", "nullable": true, "primary": false},
		{"name": "login", "type": "TEXT", "nullable": false, "primary": false},
		{"name": "password", "type": "TEXT", "nullable": false, "primary": false},
		{"name": "role", "type": "TEXT", "nullable": false, "primary": false},
		{"name": "created_at", "type": "TIMESTAMP", "nullable": false, "primary": false},
	}

	c.JSON(http.StatusOK, schema)
}
