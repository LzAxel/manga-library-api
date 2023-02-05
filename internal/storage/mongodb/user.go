package mongodb

import (
	"context"
	"errors"
	"manga-library/internal/domain"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type UserMongoDB struct {
	db *mongo.Database
}

func NewUserMongoDB(db *mongo.Database) *UserMongoDB {
	return &UserMongoDB{db: db}
}

func (m *UserMongoDB) GetByID(ctx context.Context, userID string) (domain.User, error) {
	var user domain.User

	coll := m.db.Collection(userCollection)

	cur := coll.FindOne(ctx, bson.M{"_id": userID})
	cur.Decode(&user)

	if errors.Is(cur.Err(), mongo.ErrNoDocuments) {
		return user, domain.ErrNotFound
	}

	return user, cur.Err()
}

func (m *UserMongoDB) GetByUsername(ctx context.Context, username string) (domain.User, error) {
	var user domain.User

	coll := m.db.Collection(userCollection)

	cur := coll.FindOne(ctx, bson.M{"username": username})
	cur.Decode(&user)

	if errors.Is(cur.Err(), mongo.ErrNoDocuments) {
		return user, domain.ErrNotFound
	}

	return user, cur.Err()
}
