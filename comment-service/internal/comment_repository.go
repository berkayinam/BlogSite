package internal

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type CommentRepository struct {
	collection *mongo.Collection
}

func NewCommentRepository() *CommentRepository {
	collection := Client.Database("commentdb").Collection("comments")
	return &CommentRepository{collection: collection}
}

func (r *CommentRepository) CreateComment(ctx context.Context, comment *Comment) error {
	comment.CreatedAt = time.Now()
	comment.UpdatedAt = time.Now()
	_, err := r.collection.InsertOne(ctx, comment)
	return err
}

func (r *CommentRepository) GetCommentsByPost(ctx context.Context, postID string) ([]Comment, error) {
	opts := options.Find().SetSort(bson.D{{Key: "createdAt", Value: -1}})
	cursor, err := r.collection.Find(ctx, bson.M{"postId": postID}, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var comments []Comment
	if err = cursor.All(ctx, &comments); err != nil {
		return nil, err
	}
	return comments, nil
}

func (r *CommentRepository) UpdateComment(ctx context.Context, commentID primitive.ObjectID, content string) error {
	_, err := r.collection.UpdateOne(
		ctx,
		bson.M{"_id": commentID},
		bson.M{
			"$set": bson.M{
				"content":   content,
				"updatedAt": time.Now(),
			},
		},
	)
	return err
}

func (r *CommentRepository) DeleteComment(ctx context.Context, commentID primitive.ObjectID) error {
	_, err := r.collection.DeleteOne(ctx, bson.M{"_id": commentID})
	return err
}

func (r *CommentRepository) LikeComment(ctx context.Context, commentID primitive.ObjectID, username string) error {
	_, err := r.collection.UpdateOne(
		ctx,
		bson.M{"_id": commentID},
		bson.M{"$addToSet": bson.M{"likes": username}},
	)
	return err
}

func (r *CommentRepository) UnlikeComment(ctx context.Context, commentID primitive.ObjectID, username string) error {
	_, err := r.collection.UpdateOne(
		ctx,
		bson.M{"_id": commentID},
		bson.M{"$pull": bson.M{"likes": username}},
	)
	return err
}

// Friendship repository methods
type FriendshipRepository struct {
	collection *mongo.Collection
}

func NewFriendshipRepository() *FriendshipRepository {
	collection := Client.Database("commentdb").Collection("friendships")
	return &FriendshipRepository{collection: collection}
}

func (r *FriendshipRepository) SendFriendRequest(ctx context.Context, friendship *Friendship) error {
	friendship.CreatedAt = time.Now()
	friendship.Status = "pending"
	_, err := r.collection.InsertOne(ctx, friendship)
	return err
}

func (r *FriendshipRepository) GetFriendRequests(ctx context.Context, username string) ([]Friendship, error) {
	cursor, err := r.collection.Find(ctx, bson.M{
		"user2":  username,
		"status": "pending",
	})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var requests []Friendship
	if err = cursor.All(ctx, &requests); err != nil {
		return nil, err
	}
	return requests, nil
}

func (r *FriendshipRepository) AcceptFriendRequest(ctx context.Context, friendshipID primitive.ObjectID) error {
	_, err := r.collection.UpdateOne(
		ctx,
		bson.M{"_id": friendshipID},
		bson.M{
			"$set": bson.M{
				"status":    "accepted",
				"updatedAt": time.Now(),
			},
		},
	)
	return err
}

func (r *FriendshipRepository) GetFriends(ctx context.Context, username string) ([]string, error) {
	cursor, err := r.collection.Find(ctx, bson.M{
		"$or": []bson.M{
			{"user1": username},
			{"user2": username},
		},
		"status": "accepted",
	})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var friendships []Friendship
	if err = cursor.All(ctx, &friendships); err != nil {
		return nil, err
	}

	friends := make([]string, 0)
	for _, f := range friendships {
		if f.User1 == username {
			friends = append(friends, f.User2)
		} else {
			friends = append(friends, f.User1)
		}
	}
	return friends, nil
}

// Post like repository methods
type PostLikeRepository struct {
	collection *mongo.Collection
}

func NewPostLikeRepository() *PostLikeRepository {
	collection := Client.Database("commentdb").Collection("postlikes")
	return &PostLikeRepository{collection: collection}
}

func (r *PostLikeRepository) LikePost(ctx context.Context, like *PostLike) error {
	like.CreatedAt = time.Now()
	_, err := r.collection.InsertOne(ctx, like)
	return err
}

func (r *PostLikeRepository) UnlikePost(ctx context.Context, postID string, username string) error {
	_, err := r.collection.DeleteOne(ctx, bson.M{
		"postId":   postID,
		"username": username,
	})
	return err
}

func (r *PostLikeRepository) GetPostLikes(ctx context.Context, postID string) ([]string, error) {
	cursor, err := r.collection.Find(ctx, bson.M{"postId": postID})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var likes []PostLike
	if err = cursor.All(ctx, &likes); err != nil {
		return nil, err
	}

	usernames := make([]string, len(likes))
	for i, like := range likes {
		usernames[i] = like.Username
	}
	return usernames, nil
}

func (r *PostLikeRepository) HasUserLikedPost(ctx context.Context, postID string, username string) bool {
	count, err := r.collection.CountDocuments(ctx, bson.M{
		"postId":   postID,
		"username": username,
	})
	return err == nil && count > 0
} 