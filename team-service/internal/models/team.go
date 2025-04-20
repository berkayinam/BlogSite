package models

import (
	"time"

	"gorm.io/gorm"
)

type Team struct {
	ID          uint           `gorm:"primaryKey" json:"id"`
	Name        string         `gorm:"size:255;not null;unique" json:"name"`
	Description string         `gorm:"type:text" json:"description"`
	CreatedBy   string         `gorm:"size:255;not null" json:"created_by"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`
	Members     []TeamMember   `gorm:"foreignKey:TeamID" json:"members,omitempty"`
	Invites     []TeamInvite   `gorm:"foreignKey:TeamID" json:"invites,omitempty"`
}

type TeamMember struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	TeamID    uint           `gorm:"not null" json:"team_id"`
	Username  string         `gorm:"size:255;not null" json:"username"`
	Role      string         `gorm:"size:50;not null;default:'member'" json:"role"` // admin, member
	JoinedAt  time.Time      `json:"joined_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`
}

type TeamInvite struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	TeamID    uint           `gorm:"not null" json:"team_id"`
	Username  string         `gorm:"size:255;not null" json:"username"`
	Status    string         `gorm:"size:50;not null;default:'pending'" json:"status"` // pending, accepted, rejected
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	ExpiresAt time.Time      `json:"expires_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`
}

type JoinRequest struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	TeamID    uint           `gorm:"not null" json:"team_id"`
	Username  string         `gorm:"size:255;not null" json:"username"`
	Message   string         `gorm:"type:text" json:"message"`
	Status    string         `gorm:"size:50;not null;default:'pending'" json:"status"` // pending, accepted, rejected
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`
} 