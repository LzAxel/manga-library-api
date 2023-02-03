package domain

import (
	"errors"
	"time"
)

var (
	ErrMangaTitleExists = errors.New("manga title already exists")
	ErrFailedToGet      = errors.New("failed to get manga")
)

type Manga struct {
	Id                string    `json:"_id" bson:"_id"`
	Title             string    `json:"title" bson:"title"`
	Slug              string    `json:"slug" bson:"slug"`
	AlternativeTitles []string  `json:"alternativeTitles" bson:"alternativeTitles"`
	Description       string    `json:"description" bson:"description"`
	Tags              []string  `json:"tags" bson:"tags"`
	PreviewURL        string    `json:"previewUrl" bson:"previewUrl"`
	LikeCount         int       `json:"likeCount" bson:"likeCount"`
	AgeRating         int       `json:"ageRating" bson:"ageRating"`
	UploaderId        string    `json:"uploaderId" bson:"uploaderId"`
	ReleaseYear       int       `json:"releaseYear" bson:"releaseYear"`
	CreatedAt         time.Time `json:"createdAt" bson:"createdAt"`
}

type CreateMangaDTO struct {
	Title             string   `json:"title" binding:"required"`
	AlternativeTitles []string `json:"alternativeTitles"`
	Description       string   `json:"description" binding:"required"`
	Tags              []string `json:"tags" binding:"required"`
	PreviewURL        string   `json:"previewUrl" binding:"required"`
	AgeRating         int      `json:"ageRating" binding:"required"`
	ReleaseYear       int      `json:"releaseYear" binding:"required"`
}

type UpdateMangaDTO struct {
	Id                string    `json:"-" bson:"_id"`
	Title             *string   `json:"title"`
	AlternativeTitles *[]string `json:"alternativeTitles"`
	Description       *string   `json:"description"`
	Tags              *[]string `json:"tags"`
	PreviewURL        *string   `json:"previewUrl"`
	AgeRating         *int      `json:"ageRating"`
	ReleaseYear       *int      `json:"releaseYear"`
	Slug              string    `json:"-"`
}

type GetMangaDTO struct {
	Id   string
	Slug string
}
