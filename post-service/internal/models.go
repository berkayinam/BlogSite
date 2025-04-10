package internal

import "time"

type Post struct {
	Title     string    `json:"title" bson:"title"`
	Content   string    `json:"content" bson:"content"`
	Author    string    `json:"author" bson:"author"`
	CreatedAt time.Time `json:"createdAt" bson:"createdAt"`
}
