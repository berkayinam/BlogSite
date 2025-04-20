package internal

import (
	"time"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Comment struct {
	ID        primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	PostID    string            `bson:"postId" json:"postId"`
	Content   string            `bson:"content" json:"content"`
	Author    string            `bson:"author" json:"author"`
	CreatedAt time.Time         `bson:"createdAt" json:"createdAt"`
	UpdatedAt time.Time         `bson:"updatedAt,omitempty" json:"updatedAt,omitempty"`
	Likes     []string          `bson:"likes" json:"likes"` // array of usernames who liked
}

type PostLike struct {
	ID        primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	PostID    string            `bson:"postId" json:"postId"`
	Username  string            `bson:"username" json:"username"`
	CreatedAt time.Time         `bson:"createdAt" json:"createdAt"`
}

type Friendship struct {
	ID        primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	User1     string            `bson:"user1" json:"user1"`
	User2     string            `bson:"user2" json:"user2"`
	Status    string            `bson:"status" json:"status"` // pending, accepted
	CreatedAt time.Time         `bson:"createdAt" json:"createdAt"`
	UpdatedAt time.Time         `bson:"updatedAt,omitempty" json:"updatedAt,omitempty"`
} 