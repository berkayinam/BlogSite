package internal

import (
	"time"
)

type Team struct {
	ID          int       `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Members     []TeamMember `json:"members,omitempty"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type TeamMember struct {
	TeamID    int       `json:"team_id"`
	Username  string    `json:"username"`
	Role      string    `json:"role"`
	JoinedAt  time.Time `json:"joined_at"`
}

type TeamInvite struct {
	ID              int       `json:"id"`
	TeamID          int       `json:"team_id"`
	InviterUsername string    `json:"inviter_username"`
	InviteeUsername string    `json:"invitee_username"`
	Status          string    `json:"status"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
} 