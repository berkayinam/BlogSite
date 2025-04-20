package internal

import (
	"context"
	"errors"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type TeamRepository struct {
	teams   *mongo.Collection
	invites *mongo.Collection
}

func NewTeamRepository(db *mongo.Database) *TeamRepository {
	return &TeamRepository{
		teams:   db.Collection("teams"),
		invites: db.Collection("team_invites"),
	}
}

func (r *TeamRepository) CreateTeam(ctx context.Context, team *Team) error {
	team.CreatedAt = time.Now()
	team.UpdatedAt = time.Now()
	_, err := r.teams.InsertOne(ctx, team)
	return err
}

func (r *TeamRepository) GetTeam(ctx context.Context, id primitive.ObjectID) (*Team, error) {
	var team Team
	err := r.teams.FindOne(ctx, bson.M{"_id": id}).Decode(&team)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, errors.New("team not found")
		}
		return nil, err
	}
	return &team, nil
}

func (r *TeamRepository) GetUserTeams(ctx context.Context, username string) ([]Team, error) {
	cursor, err := r.teams.Find(ctx, bson.M{"members.username": username})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var teams []Team
	if err = cursor.All(ctx, &teams); err != nil {
		return nil, err
	}
	return teams, nil
}

func (r *TeamRepository) UpdateTeam(ctx context.Context, team *Team) error {
	team.UpdatedAt = time.Now()
	_, err := r.teams.ReplaceOne(ctx, bson.M{"_id": team.ID}, team)
	return err
}

func (r *TeamRepository) DeleteTeam(ctx context.Context, id primitive.ObjectID) error {
	_, err := r.teams.DeleteOne(ctx, bson.M{"_id": id})
	return err
}

func (r *TeamRepository) AddMember(ctx context.Context, teamID primitive.ObjectID, member TeamMember) error {
	member.JoinedAt = time.Now()
	_, err := r.teams.UpdateOne(
		ctx,
		bson.M{"_id": teamID},
		bson.M{
			"$push": bson.M{"members": member},
			"$set":  bson.M{"updatedAt": time.Now()},
		},
	)
	return err
}

func (r *TeamRepository) RemoveMember(ctx context.Context, teamID primitive.ObjectID, username string) error {
	_, err := r.teams.UpdateOne(
		ctx,
		bson.M{"_id": teamID},
		bson.M{
			"$pull": bson.M{"members": bson.M{"username": username}},
			"$set":  bson.M{"updatedAt": time.Now()},
		},
	)
	return err
}

func (r *TeamRepository) CreateInvite(ctx context.Context, invite *TeamInvite) error {
	invite.CreatedAt = time.Now()
	invite.Status = "pending"
	_, err := r.invites.InsertOne(ctx, invite)
	return err
}

func (r *TeamRepository) GetInvite(ctx context.Context, id primitive.ObjectID) (*TeamInvite, error) {
	var invite TeamInvite
	err := r.invites.FindOne(ctx, bson.M{"_id": id}).Decode(&invite)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, errors.New("invite not found")
		}
		return nil, err
	}
	return &invite, nil
}

func (r *TeamRepository) UpdateInviteStatus(ctx context.Context, id primitive.ObjectID, status string) error {
	_, err := r.invites.UpdateOne(
		ctx,
		bson.M{"_id": id},
		bson.M{
			"$set": bson.M{
				"status":    status,
				"updatedAt": time.Now(),
			},
		},
	)
	return err
}

func (r *TeamRepository) GetPendingInvites(ctx context.Context, username string) ([]TeamInvite, error) {
	cursor, err := r.invites.Find(ctx, bson.M{
		"username": username,
		"status":   "pending",
	})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var invites []TeamInvite
	if err = cursor.All(ctx, &invites); err != nil {
		return nil, err
	}
	return invites, nil
} 