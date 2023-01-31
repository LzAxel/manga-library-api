package mongodb

import (
	"context"
	"manga-library/internal/domain"
	"manga-library/pkg/errors"
	"manga-library/pkg/logger"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

const (
	userCollection = "user"
)

type AuthorizationMongoDB struct {
	logger logger.Logger
	db     *mongo.Database
}

func NewAuthorizationMongoDB(logger logger.Logger, db *mongo.Database) *AuthorizationMongoDB {
	return &AuthorizationMongoDB{
		logger: logger,
		db:     db,
	}
}

func (s *AuthorizationMongoDB) SignUp(ctx context.Context, user domain.User) error {
	coll := s.db.Collection(userCollection)

	result, err := coll.InsertOne(ctx, user)
	if err != nil {
		if mongo.IsDuplicateKeyError(err) {
			return errors.ErrUsernameExists
		}
		return err
	}
	s.logger.Debugf("insertedID: %v", result.InsertedID)

	return nil
}

func (s *AuthorizationMongoDB) SignIn(ctx context.Context, username string) (password, userId string, err error) {
	var user domain.User

	coll := s.db.Collection(userCollection)

	cur := coll.FindOne(ctx, bson.M{"username": username})
	cur.Decode(&user)

	return user.PasswordHash, user.Id, cur.Err()
}
