package domain

import (
	"mime/multipart"
	"time"
)

type Chapter struct {
	Volume      int    `json:"volume" bson:"volume"`
	Number      int    `json:"number" bson:"number"`
	PageCount   int    `json:"pageCount" bson:"pageCount"`
	MangaSlug   string `json:"-" bson:"-"`
	UploaderId  string `json:"uploaderId" bson:"uploaderId"`
	IsPublished bool   `json:"isPublished" bson:"isPublished"`

	UploadedAt time.Time `json:"uploadedAt" bson:"uploadedAt"`
}

type UploadChapterDTO struct {
	MangaSlug  string                `form:"mangaSlug" binding:"required"`
	Volume     int                   `form:"volume" binding:"required,numeric,gte=0"`
	Number     int                   `form:"number" binding:"required,numeric,gte=0"`
	UploaderID string                `form:"-"`
	File       *multipart.FileHeader `form:"file" binding:"required"`
}
