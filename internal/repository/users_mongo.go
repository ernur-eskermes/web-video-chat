package repository

import (
	"context"
	"errors"
	"time"

	"github.com/ernur-eskermes/web-video-chat/internal/core"
	"github.com/ernur-eskermes/web-video-chat/pkg/database/mongodb"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type UsersRepo struct {
	db *mongo.Collection
}

func NewUsersRepo(db *mongo.Database) *UsersRepo {
	return &UsersRepo{
		db: db.Collection(usersCollection),
	}
}

func (r *UsersRepo) GetById(ctx context.Context, id primitive.ObjectID) (core.User, error) {
	var user core.User

	if err := r.db.FindOne(ctx, bson.M{"_id": id}).Decode(&user); err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return core.User{}, core.ErrUserNotFound
		}

		return core.User{}, err
	}

	return user, nil
}

func (r *UsersRepo) Create(ctx context.Context, user *core.User) error {
	res, err := r.db.InsertOne(ctx, user)
	if mongodb.IsDuplicate(err) {
		return core.ErrUserAlreadyExists
	}

	user.ID = res.InsertedID.(primitive.ObjectID) //nolint:forcetypeassert

	return err
}

func (r *UsersRepo) CreateSubscription(ctx context.Context, userId, subscription primitive.ObjectID) error {
	_, err := r.db.UpdateOne(ctx, bson.M{"_id": userId}, bson.M{"$addToSet": bson.M{"pnd_subs": subscription}})
	if err != nil {
		return err
	}

	return nil
}

func (r *UsersRepo) GetByCredentials(ctx context.Context, username, password, provider string) (core.User, error) {
	var user core.User
	if err := r.db.FindOne(ctx, bson.M{"username": username, "password": password, "provider": provider}).Decode(&user); err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return core.User{}, core.ErrUserNotFound
		}

		return core.User{}, err
	}

	return user, nil
}

func (r *UsersRepo) GetByRefreshToken(ctx context.Context, refreshToken string) (core.User, error) {
	var user core.User
	if err := r.db.FindOne(ctx, bson.M{
		"session.refreshToken": refreshToken,
		"session.expiresAt":    bson.M{"$gt": time.Now()},
	}).Decode(&user); err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return core.User{}, core.ErrUserNotFound
		}

		return core.User{}, err
	}

	return user, nil
}

func (r *UsersRepo) SetSession(ctx context.Context, userID primitive.ObjectID, session core.Session) error {
	_, err := r.db.UpdateOne(ctx, bson.M{"_id": userID}, bson.M{"$set": bson.M{"session": session, "lastVisitAt": time.Now()}})

	return err
}
