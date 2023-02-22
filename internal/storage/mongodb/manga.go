package mongodb

import (
	"context"
	"errors"
	"log"
	"manga-library/internal/domain"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	mangaCollection  = "manga"
	latestMangaLimit = 20
)

type MangaMongoDB struct {
	db *mongo.Database
}

func NewMangaMongoDB(db *mongo.Database) *MangaMongoDB {
	return &MangaMongoDB{db: db}
}

func (m *MangaMongoDB) Create(ctx context.Context, manga domain.Manga) (string, error) {
	coll := m.db.Collection(mangaCollection)

	result, err := coll.InsertOne(ctx, manga)
	if err != nil {
		return "", err
	}

	id, ok := result.InsertedID.(string)
	if !ok {
		return "", errors.New("failed to get insertedId")
	}

	return id, nil
}

func (m *MangaMongoDB) GetLatest(ctx context.Context) ([]domain.Manga, error) {
	var result []domain.Manga

	coll := m.db.Collection(mangaCollection)
	opts := options.Find().SetLimit(latestMangaLimit).SetSort(bson.D{{Key: "createdAt", Value: -1}})

	cur, err := coll.Find(ctx, bson.D{{Key: "isPublished", Value: true}}, opts)
	if err != nil {
		return []domain.Manga{}, err
	}
	if err := cur.All(ctx, &result); err != nil {
		return []domain.Manga{}, err
	}

	return result, nil
}

func (m *MangaMongoDB) GetById(ctx context.Context, mangaId string) (domain.Manga, error) {
	var manga domain.Manga

	coll := m.db.Collection(mangaCollection)

	cur := coll.FindOne(ctx, bson.M{"_id": mangaId})
	cur.Decode(&manga)

	return manga, cur.Err()
}

func (m *MangaMongoDB) GetBySlug(ctx context.Context, mangaSlug string) (domain.Manga, error) {
	var manga domain.Manga

	coll := m.db.Collection(mangaCollection)

	cur := coll.FindOne(ctx, bson.M{"slug": mangaSlug})
	cur.Decode(&manga)

	return manga, cur.Err()
}

func (m *MangaMongoDB) GetByTags(ctx context.Context, tags []string, offset int) ([]domain.Manga, error) {
	var result []domain.Manga

	coll := m.db.Collection(mangaCollection)
	filter := bson.D{
		{Key: "isPublished", Value: true},
		{Key: "tags", Value: bson.D{
			{Key: "$all", Value: tags},
		}},
	}

	opts := options.Find().SetLimit(latestMangaLimit).
		SetSort(bson.D{{Key: "createdAt", Value: -1}}).
		SetSkip(int64(offset))

	cur, err := coll.Find(ctx, filter, opts)
	if err != nil {
		return []domain.Manga{}, err
	}
	if err := cur.All(ctx, &result); err != nil {
		return []domain.Manga{}, err
	}

	return result, nil
}

func (m *MangaMongoDB) Delete(ctx context.Context, mangaId string) error {

	coll := m.db.Collection(mangaCollection)

	if _, err := coll.DeleteOne(ctx, bson.M{"_id": mangaId}); err != nil {
		return err
	}

	return nil
}

func (m *MangaMongoDB) Update(ctx context.Context, mangaDTO domain.UpdateMangaDTO) error {
	coll := m.db.Collection(mangaCollection)

	var setQuery bson.D

	if mangaDTO.Title != nil {
		setQuery = append(setQuery, bson.E{Key: "title", Value: &mangaDTO.Title})
		setQuery = append(setQuery, bson.E{Key: "slug", Value: mangaDTO.Slug})
	}
	if mangaDTO.AlternativeTitles != nil {
		setQuery = append(setQuery, bson.E{Key: "alternativeTitles", Value: &mangaDTO.AlternativeTitles})
	}
	if mangaDTO.Description != nil {
		setQuery = append(setQuery, bson.E{Key: "description", Value: &mangaDTO.Description})
	}
	if mangaDTO.Tags != nil {
		setQuery = append(setQuery, bson.E{Key: "tags", Value: &mangaDTO.Tags})
	}
	if mangaDTO.AgeRating != nil {
		setQuery = append(setQuery, bson.E{Key: "ageRating", Value: &mangaDTO.AgeRating})
	}
	if mangaDTO.ReleaseYear != nil {
		setQuery = append(setQuery, bson.E{Key: "releaseYear", Value: &mangaDTO.ReleaseYear})
	}
	if mangaDTO.IsPublished != nil {
		setQuery = append(setQuery, bson.E{Key: "isPublished", Value: &mangaDTO.IsPublished})
	}
	if mangaDTO.Author != nil {
		setQuery = append(setQuery, bson.E{Key: "author", Value: &mangaDTO.Author})
	}

	if setQuery == nil {
		return nil
	}
	result, err := coll.UpdateByID(ctx, mangaDTO.ID, bson.D{{Key: "$set", Value: setQuery}})
	if err != nil {
		return err
	}
	log.Println(result)

	return nil
}
