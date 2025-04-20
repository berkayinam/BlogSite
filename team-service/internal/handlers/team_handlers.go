package handlers

import (
	"net/http"
	"time"

	"github.com/binam/myblog-project/team-service/internal/database"
	"github.com/binam/myblog-project/team-service/internal/models"
	"github.com/gin-gonic/gin"
)

// CreateTeam handles team creation
func CreateTeam(c *gin.Context) {
	var team models.Team
	if err := c.ShouldBindJSON(&team); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Get username from JWT token
	username := c.GetString("username")
	if username == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	team.CreatedBy = username
	team.CreatedAt = time.Now()
	team.UpdatedAt = time.Now()

	if err := database.DB.Create(&team).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Add creator as team admin
	member := models.TeamMember{
		TeamID:   team.ID,
		Username: username,
		Role:     "admin",
		JoinedAt: time.Now(),
	}

	if err := database.DB.Create(&member).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, team)
}

// InviteToTeam handles team invitations
func InviteToTeam(c *gin.Context) {
	var invite models.TeamInvite
	if err := c.ShouldBindJSON(&invite); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Check if user is team admin
	username := c.GetString("username")
	var member models.TeamMember
	if err := database.DB.Where("team_id = ? AND username = ? AND role = ?", 
		invite.TeamID, username, "admin").First(&member).Error; err != nil {
		c.JSON(http.StatusForbidden, gin.H{"error": "only team admins can send invites"})
		return
	}

	invite.CreatedAt = time.Now()
	invite.UpdatedAt = time.Now()
	invite.ExpiresAt = time.Now().Add(7 * 24 * time.Hour) // Expires in 7 days
	invite.Status = "pending"

	if err := database.DB.Create(&invite).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, invite)
}

// RespondToInvite handles responses to team invitations
func RespondToInvite(c *gin.Context) {
	var response struct {
		InviteID uint   `json:"invite_id"`
		Status   string `json:"status"` // accepted or rejected
	}

	if err := c.ShouldBindJSON(&response); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	username := c.GetString("username")
	var invite models.TeamInvite
	if err := database.DB.Where("id = ? AND username = ?", 
		response.InviteID, username).First(&invite).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "invite not found"})
		return
	}

	if invite.Status != "pending" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invite already processed"})
		return
	}

	if time.Now().After(invite.ExpiresAt) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invite expired"})
		return
	}

	invite.Status = response.Status
	invite.UpdatedAt = time.Now()

	if err := database.DB.Save(&invite).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if response.Status == "accepted" {
		member := models.TeamMember{
			TeamID:   invite.TeamID,
			Username: username,
			Role:     "member",
			JoinedAt: time.Now(),
		}

		if err := database.DB.Create(&member).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
	}

	c.JSON(http.StatusOK, gin.H{"message": "invite processed successfully"})
}

// RequestToJoinTeam handles join requests
func RequestToJoinTeam(c *gin.Context) {
	var request models.JoinRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	username := c.GetString("username")
	request.Username = username
	request.Status = "pending"
	request.CreatedAt = time.Now()
	request.UpdatedAt = time.Now()

	if err := database.DB.Create(&request).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, request)
}

// GetTeam retrieves team details
func GetTeam(c *gin.Context) {
	teamID := c.Param("id")
	var team models.Team

	if err := database.DB.Preload("Members").First(&team, teamID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "team not found"})
		return
	}

	c.JSON(http.StatusOK, team)
}

// ListTeams lists all teams (with optional filters)
func ListTeams(c *gin.Context) {
	var teams []models.Team
	query := database.DB.Preload("Members")

	// Add filters if needed
	if search := c.Query("search"); search != "" {
		query = query.Where("name LIKE ?", "%"+search+"%")
	}

	if err := query.Find(&teams).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, teams)
} 