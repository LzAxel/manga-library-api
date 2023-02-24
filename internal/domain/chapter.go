package domain

import (
	"mime/multipart"
	"time"
)

type Chapter struct {
	Volume      int    `json:"volume" bson:"volume"`
	Number      int    `json:"number" bson:"number"`
	PageCount   int    `json:"pageCount" bson:"pageCount"`
	Title       string `json:"title" bson:"title"`
	MangaSlug   string `json:"-" bson:"-"`
	UploaderId  string `json:"uploaderId" bson:"uploaderId"`
	IsPublished bool   `json:"isPublished" bson:"isPublished"`

	UploadedAt time.Time `json:"uploadedAt" bson:"uploadedAt"`
}

type UploadChapterDTO struct {
	Title      string                `form:"title" binding:"required"`
	MangaSlug  string                `form:"mangaSlug" binding:"required"`
	UploaderID string                `form:"-" swaggerignore:"true"`
	Volume     int                   `form:"volume" binding:"required,numeric,gte=0"`
	Number     int                   `form:"number" binding:"required,numeric,gte=0"`
	File       *multipart.FileHeader `form:"file" binding:"required" swaggerignore:"true"`
}

type DeleteChapterDTO struct {
	MangaSlug  string
	UploaderID string
	Volume     int
	Number     int
}
