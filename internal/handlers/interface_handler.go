package handlers

import (
	"net/http"
	"strconv"

	"gosql_interface/internal/database"
	"gosql_interface/internal/models"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

func ListInterfacesPage(c *gin.Context) {
	session := sessions.Default(c)
	c.HTML(http.StatusOK, "interfaces.html", gin.H{
		"title":      "Interfaces",
		"activePage": "interfaces",
		"userLogin":  session.Get("user_login"),
		"userRole":   session.Get("user_role"),
	})
}

func GetInterfacesAPI(c *gin.Context) {
	search := c.DefaultQuery("search", "")
	db := database.GetDB()

	interfaces, err := models.GetAllInterfaces(db, search)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, interfaces)
}

func CreateInterfaceAPI(c *gin.Context) {
	var iface models.Interface
	if err := c.ShouldBindJSON(&iface); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid data"})
		return
	}

	db := database.GetDB()
	if err := models.CreateInterface(db, &iface); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, iface)
}

func UpdateInterfaceAPI(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	var iface models.Interface
	if err := c.ShouldBindJSON(&iface); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid data"})
		return
	}

	iface.ID = id
	db := database.GetDB()
	if err := models.UpdateInterface(db, &iface); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, iface)
}

func DeleteInterfaceAPI(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	db := database.GetDB()
	if err := models.DeleteInterface(db, id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Interface deleted successfully"})
}

func GetInterfaceSchemaAPI(c *gin.Context) {
	schema := []map[string]interface{}{
		{"name": "id", "type": "BIGINT", "nullable": false, "primary": true},
		{"name": "name", "type": "TEXT", "nullable": false, "primary": false},
		{"name": "ip", "type": "TEXT", "nullable": false, "primary": false},
		{"name": "mac", "type": "TEXT", "nullable": false, "primary": false},
		{"name": "mtu", "type": "INTEGER", "nullable": false, "primary": false},
		{"name": "status", "type": "BOOLEAN", "nullable": false, "primary": false},
		{"name": "ip_type", "type": "BOOLEAN", "nullable": false, "primary": false},
		{"name": "created_at", "type": "TIMESTAMP", "nullable": false, "primary": false},
	}

	c.JSON(http.StatusOK, schema)
}
