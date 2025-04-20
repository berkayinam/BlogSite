package internal

import (
	"time"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Team struct {
	ID          primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	Name        string            `bson:"name" json:"name"`
	Description string            `bson:"description" json:"description"`
	Members     []TeamMember      `bson:"members" json:"members"`
	CreatedAt   time.Time         `bson:"createdAt" json:"createdAt"`
	UpdatedAt   time.Time         `bson:"updatedAt,omitempty" json:"updatedAt,omitempty"`
}

type TeamMember struct {
	Username  string    `bson:"username" json:"username"`
	Role      string    `bson:"role" json:"role"` // admin or member
	JoinedAt  time.Time `bson:"joinedAt" json:"joinedAt"`
}

type TeamInvite struct {
	ID        primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	TeamID    primitive.ObjectID `bson:"teamId" json:"teamId"`
	Username  string            `bson:"username" json:"username"`
	Status    string            `bson:"status" json:"status"` // pending, accepted, rejected
	CreatedAt time.Time         `bson:"createdAt" json:"createdAt"`
	UpdatedAt time.Time         `bson:"updatedAt,omitempty" json:"updatedAt,omitempty"`
} 