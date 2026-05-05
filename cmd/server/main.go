package main

import (
	"fmt"
	"log"
	"net/http"

	"gosql_interface/internal/config"
	"gosql_interface/internal/database"
	"gosql_interface/internal/handlers"
	"gosql_interface/internal/middleware"
	"gosql_interface/internal/models"

	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
)

func main() {
	cfg := config.LoadConfig()

	if err := database.InitDB(cfg.DatabaseURL); err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}
	defer database.CloseDB()

	if err := database.RunMigrations(); err != nil {
		log.Fatalf("Failed to run migrations: %v", err)
	}

	if cfg.Environment == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	router := gin.Default()

	store := cookie.NewStore([]byte(cfg.SessionSecret))
	router.Use(sessions.Sessions("gosql_session", store))

	router.Static("/static", "./web/static")

	router.LoadHTMLGlob("web/templates/**/*")

	router.GET("/login", handlers.LoginPage)
	router.POST("/login", handlers.LoginHandler)
	router.GET("/logout", handlers.LogoutHandler)

	router.GET("/", func(c *gin.Context) {
		c.Redirect(http.StatusFound, "/interfaces")
	})

	authorized := router.Group("/")
	authorized.Use(middleware.AuthRequired())
	{
		authorized.GET("/interfaces", handlers.ListInterfacesPage)
		authorized.GET("/api/interfaces", handlers.GetInterfacesAPI)
		authorized.GET("/api/interfaces/schema", handlers.GetInterfaceSchemaAPI)
		authorized.GET("/api/interfaces/:id", getInterfaceByIDHandler)

		authorized.GET("/acl", handlers.ListACLsPage)
		authorized.GET("/api/acl", handlers.GetACLsAPI)
		authorized.GET("/api/acl/schema", handlers.GetACLSchemaAPI)
		authorized.GET("/api/acl/:id", getACLByIDHandler)
	}

	admin := router.Group("/")
	admin.Use(middleware.AuthRequired(), middleware.AdminRequired())
	{
		admin.GET("/users", handlers.ListUsersPage)
		admin.GET("/api/users", handlers.GetUsersAPI)
		admin.GET("/api/users/schema", handlers.GetUserSchemaAPI)
		admin.GET("/api/users/:id", getUserByIDHandler)
		admin.POST("/api/users", handlers.CreateUserAPI)
		admin.PUT("/api/users/:id", handlers.UpdateUserAPI)
		admin.DELETE("/api/users/:id", handlers.DeleteUserAPI)

		admin.POST("/api/interfaces", handlers.CreateInterfaceAPI)
		admin.PUT("/api/interfaces/:id", handlers.UpdateInterfaceAPI)
		admin.DELETE("/api/interfaces/:id", handlers.DeleteInterfaceAPI)

		admin.POST("/api/acl", handlers.CreateACLAPI)
		admin.PUT("/api/acl/:id", handlers.UpdateACLAPI)
		admin.DELETE("/api/acl/:id", handlers.DeleteACLAPI)
	}

	log.Printf("✓ Server starting on port %s", cfg.ServerPort)
	log.Printf("✓ Open http://localhost:%s/login in your browser", cfg.ServerPort)
	log.Printf("✓ Default credentials: admin / admin123")

	if err := router.Run(":" + cfg.ServerPort); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}

func getUserByIDHandler(c *gin.Context) {
	id := c.Param("id")
	db := database.GetDB()

	var idInt int64
	if _, err := fmt.Sscanf(id, "%d", &idInt); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	user, err := models.GetUserByID(db, idInt)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	user.Password = ""
	c.JSON(http.StatusOK, user)
}

func getInterfaceByIDHandler(c *gin.Context) {
	id := c.Param("id")
	db := database.GetDB()

	var idInt int64
	if _, err := fmt.Sscanf(id, "%d", &idInt); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	iface, err := models.GetInterfaceByID(db, idInt)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, iface)
}

func getACLByIDHandler(c *gin.Context) {
	id := c.Param("id")
	db := database.GetDB()

	var idInt int64
	if _, err := fmt.Sscanf(id, "%d", &idInt); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	acl, err := models.GetACLByID(db, idInt)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, acl)
}
