package main

import (
	"github.com/gin-gonic/gin"
	"team-service/internal"
)

func main() {
	r := gin.Default()

	// Initialize MongoDB connection
	if err := internal.ConnectToMongo(); err != nil {
		panic(err)
	}

	// Initialize repositories and handlers
	teamRepo := internal.NewTeamRepository(internal.Client.Database("myblog"))
	handler := internal.NewTeamHandler(teamRepo)

	// Routes
	r.POST("/teams", internal.AuthMiddleware(handler.CreateTeam))
	r.GET("/teams/:id", handler.GetTeam)
	r.GET("/teams/user", internal.AuthMiddleware(handler.GetUserTeams))
	r.PUT("/teams/:id", internal.AuthMiddleware(handler.UpdateTeam))
	r.DELETE("/teams/:id", internal.AuthMiddleware(handler.DeleteTeam))
	r.POST("/teams/:id/invite", internal.AuthMiddleware(handler.InviteMember))
	r.POST("/teams/invites/:id", internal.AuthMiddleware(handler.RespondToInvite))
	r.DELETE("/teams/:id/members/:username", internal.AuthMiddleware(handler.RemoveMember))

	// Health check
	r.GET("/health", func(c *gin.Context) {
		c.Status(200)
	})

	r.Run(":8084")
} 