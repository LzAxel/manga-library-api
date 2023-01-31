package mongodb

import (
	"context"
	"manga-library/internal/domain"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

const previewCollection = "preview"

type PreviewMongoDB struct {
	db *mongo.Database
}

func NewPreviewMongoDB(db *mongo.Database) *PreviewMongoDB {
	return &PreviewMongoDB{db: db}
}

func (m *PreviewMongoDB) Create(ctx context.Context, preview domain.Preview) (string, error) {
	coll := m.db.Collection(previewCollection)

	_, err := coll.InsertOne(ctx, preview)
	if err != nil {
		return "", err
	}

	return preview.URL, nil
}

func (m *PreviewMongoDB) Delete(ctx context.Context, previewId string) error {
	coll := m.db.Collection(previewCollection)

	_, err := coll.DeleteOne(ctx, bson.M{"_id": previewId})
	if err != nil {
		return err
	}

	return nil
}
