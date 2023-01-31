package domain

import (
	"time"
)

type Preview struct {
	Id         string    `json:"_id" bson:"_id"`
	FileName   string    `json:"fileName" bson:"fileName"`
	UploaderId string    `json:"uploaderId" bson:"uploaderId"`
	URL        string    `json:"url" bson:"url"`
	CreatedAt  time.Time `json:"createdAt" bson:"createdAt"`
}

type CreatePreviewDTO struct {
	UploaderId string    `json:"uploaderId"`
	URL        string    `json:"url"`
	CreatedAt  time.Time `json:"createdAt"`
}
