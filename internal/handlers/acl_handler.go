package handlers

import (
	"net/http"
	"strconv"

	"gosql_interface/internal/database"
	"gosql_interface/internal/models"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

// ListACLsPage renders the ACL management page
func ListACLsPage(c *gin.Context) {
	session := sessions.Default(c)
	c.HTML(http.StatusOK, "acl.html", gin.H{
		"title":      "ACL",
		"activePage": "acl",
		"userLogin":  session.Get("user_login"),
	})
}

// GetACLsAPI returns ACL rules as JSON with optional search
func GetACLsAPI(c *gin.Context) {
	search := c.DefaultQuery("search", "")
	db := database.GetDB()

	acls, err := models.GetAllACLs(db, search)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, acls)
}

// CreateACLAPI creates a new ACL rule
func CreateACLAPI(c *gin.Context) {
	var acl models.ACL
	if err := c.ShouldBindJSON(&acl); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid data"})
		return
	}

	db := database.GetDB()
	if err := models.CreateACL(db, &acl); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, acl)
}

// UpdateACLAPI updates an existing ACL rule
func UpdateACLAPI(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	var acl models.ACL
	if err := c.ShouldBindJSON(&acl); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid data"})
		return
	}

	acl.ID = id
	db := database.GetDB()
	if err := models.UpdateACL(db, &acl); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, acl)
}

// DeleteACLAPI deletes an ACL rule
func DeleteACLAPI(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	db := database.GetDB()
	if err := models.DeleteACL(db, id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "ACL rule deleted successfully"})
}

// GetACLSchemaAPI returns the schema for ACL table
func GetACLSchemaAPI(c *gin.Context) {
	schema := []map[string]interface{}{
		{"name": "id", "type": "BIGINT", "nullable": false, "primary": true},
		{"name": "src_ip", "type": "TEXT", "nullable": false, "primary": false},
		{"name": "dst_ip", "type": "TEXT", "nullable": false, "primary": false},
		{"name": "protocol", "type": "TEXT", "nullable": false, "primary": false},
		{"name": "src_port", "type": "TEXT", "nullable": true, "primary": false},
		{"name": "dst_port", "type": "TEXT", "nullable": true, "primary": false},
		{"name": "action", "type": "TEXT", "nullable": false, "primary": false},
		{"name": "created_at", "type": "TIMESTAMP", "nullable": false, "primary": false},
	}

	c.JSON(http.StatusOK, schema)
}
