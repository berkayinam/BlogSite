package main

import (
	"log"

	"github.com/binam/myblog-project/team-service/internal/database"
	"github.com/binam/myblog-project/team-service/internal/handlers"
	"github.com/binam/myblog-project/team-service/internal/middleware"
	"github.com/gin-gonic/gin"
)

func main() {
	// Initialize database connection
	database.InitDB()

	// Create Gin router
	r := gin.Default()

	// Add CORS middleware
	r.Use(func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Origin, Authorization, Content-Type")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	})

	// Public routes
	r.GET("/teams", handlers.ListTeams)
	r.GET("/teams/:id", handlers.GetTeam)

	// Protected routes
	protected := r.Group("/")
	protected.Use(middleware.AuthMiddleware())
	{
		// Team management
		protected.POST("/teams", handlers.CreateTeam)
		
		// Team invitations
		protected.POST("/teams/invite", handlers.InviteToTeam)
		protected.POST("/teams/invite/respond", handlers.RespondToInvite)
		
		// Join requests
		protected.POST("/teams/join/request", handlers.RequestToJoinTeam)
	}

	// Start server
	if err := r.Run(":8082"); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
} 