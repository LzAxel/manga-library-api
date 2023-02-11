package mongodb

import (
	"context"
	"manga-library/internal/domain"

	"go.mongodb.org/mongo-driver/bson"
)

func (m *MangaMongoDB) UploadChapter(ctx context.Context, chapter domain.Chapter) error {
	coll := m.db.Collection(mangaCollection)

	_, err := coll.UpdateOne(ctx, bson.M{"slug": chapter.MangaSlug}, bson.M{"$push": bson.M{"chapters": chapter}})
	if err != nil {
		return err
	}

	return nil
}

func (m *MangaMongoDB) DeleteChapter(ctx context.Context, chapter domain.DeleteChapterDTO) error {
	coll := m.db.Collection(mangaCollection)

	result, err := coll.UpdateOne(ctx, bson.M{"slug": chapter.MangaSlug},
		bson.M{"$pull": bson.M{"chapters": bson.M{"volume": chapter.Volume, "number": chapter.Number}}})
	if err != nil {
		return err
	}
	if result.ModifiedCount == 0 {
		return domain.ErrNotFound
	}

	return nil
}
