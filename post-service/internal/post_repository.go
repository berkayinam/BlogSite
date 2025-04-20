package internal

import (
	"context"
	"errors"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	databaseName   = "blogdb"
	collectionName = "posts"
)

type PostRepository struct {
	collection *mongo.Collection
}

func NewPostRepository() *PostRepository {
	collection := Client.Database(databaseName).Collection(collectionName)
	return &PostRepository{
		collection: collection,
	}
}

// CreatePost creates a new post in the database
func (r *PostRepository) CreatePost(ctx context.Context, post *Post) error {
	post.CreatedAt = time.Now()
	_, err := r.collection.InsertOne(ctx, post)
	if err != nil {
		return err
	}
	return nil
}

// GetAllPosts retrieves all posts from the database
func (r *PostRepository) GetAllPosts(ctx context.Context) ([]Post, error) {
	opts := options.Find().SetSort(bson.D{{Key: "createdAt", Value: -1}})
	cursor, err := r.collection.Find(ctx, bson.M{}, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var posts []Post
	if err = cursor.All(ctx, &posts); err != nil {
		return nil, err
	}
	return posts, nil
}

// GetPostsByAuthor retrieves all posts by a specific author
func (r *PostRepository) GetPostsByAuthor(ctx context.Context, author string) ([]Post, error) {
	cursor, err := r.collection.Find(ctx, bson.M{"author": author})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var posts []Post
	if err = cursor.All(ctx, &posts); err != nil {
		return nil, err
	}
	return posts, nil
}

// UpdatePost updates an existing post
func (r *PostRepository) UpdatePost(ctx context.Context, title string, post *Post) error {
	result, err := r.collection.UpdateOne(
		ctx,
		bson.M{"title": title},
		bson.M{"$set": post},
	)
	if err != nil {
		return err
	}
	if result.MatchedCount == 0 {
		return errors.New("post not found")
	}
	return nil
}

// DeletePost deletes a post by its title
func (r *PostRepository) DeletePost(ctx context.Context, title string) error {
	result, err := r.collection.DeleteOne(ctx, bson.M{"title": title})
	if err != nil {
		return err
	}
	if result.DeletedCount == 0 {
		return errors.New("post not found")
	}
	return nil
}

// SearchPosts searches for posts based on title or content
func (r *PostRepository) SearchPosts(ctx context.Context, query string) ([]Post, error) {
	filter := bson.M{
		"$or": []bson.M{
			{"title": bson.M{"$regex": primitive.Regex{Pattern: query, Options: "i"}}},
			{"content": bson.M{"$regex": primitive.Regex{Pattern: query, Options: "i"}}},
		},
	}
	
	cursor, err := r.collection.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var posts []Post
	if err = cursor.All(ctx, &posts); err != nil {
		return nil, err
	}
	return posts, nil
} 