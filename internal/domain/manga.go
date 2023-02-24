package domain

import (
	"errors"
	"fmt"
	"mime/multipart"
	"strings"
	"time"
)

const (
	orderByValues = "date releaseDate"
)

var (
	ErrMangaTitleExists   = errors.New("manga title already exists")
	ErrFailedToGet        = errors.New("failed to get manga")
	ErrMangaPublishByUser = errors.New("only admin user can publish manga")
	ErrInvalidFilter      = errors.New("invalid filter")
)

type Manga struct {
	ID                string    `json:"_id" bson:"_id"`
	Title             string    `json:"title" bson:"title"`
	Slug              string    `json:"slug" bson:"slug"`
	Author            string    `json:"author" bson:"author"`
	Chapters          []Chapter `json:"chapters" bson:"chapters"`
	Description       string    `json:"description" bson:"description"`
	PreviewURL        string    `json:"previewUrl" bson:"previewUrl"`
	UploaderId        string    `json:"uploaderId" bson:"uploaderId"`
	AlternativeTitles []string  `json:"alternativeTitles" bson:"alternativeTitles"`
	Tags              []string  `json:"tags" bson:"tags"`
	LikeCount         int       `json:"likeCount" bson:"likeCount"`
	AgeRating         int       `json:"ageRating" bson:"ageRating"`
	ReleaseYear       int       `json:"releaseYear" bson:"releaseYear"`
	IsPublished       bool      `json:"isPublished" bson:"isPublished"`
	CreatedAt         time.Time `json:"createdAt" bson:"createdAt"`
}

type CreateMangaRawDTO struct {
	Title             string                `form:"title" binding:"required"`
	AlternativeTitles string                `form:"alternativeTitles" example:"title1,title2,title3 (separate by comma)"`
	Author            string                `form:"author" binding:"required"`
	Description       string                `form:"description" binding:"required"`
	Tags              string                `form:"tags" binding:"required" example:"tag1,tag2,tag3 (separate by comma)"`
	Preview           *multipart.FileHeader `form:"preview" binding:"required" swaggerignore:"true"`
	AgeRating         int                   `form:"ageRating" binding:"required"`
	ReleaseYear       int                   `form:"releaseYear" binding:"required"`
}

func (m *CreateMangaRawDTO) Validate() CreateMangaDTO {
	var tags = make([]string, 0)
	var altTitles = make([]string, 0)

	for _, tag := range strings.Split(m.Tags, ",") {
		if len(tag) > 0 {
			tags = append(tags, tag)
		}
	}

	for _, altTitle := range strings.Split(m.AlternativeTitles, ",") {
		if len(altTitle) > 0 {
			altTitles = append(altTitles, altTitle)
		}
	}

	return CreateMangaDTO{
		Title:             m.Title,
		AlternativeTitles: altTitles,
		Author:            m.Author,
		Description:       m.Description,
		Tags:              tags,
		Preview:           m.Preview,
		AgeRating:         m.AgeRating,
		ReleaseYear:       m.ReleaseYear,
	}
}

type CreateMangaDTO struct {
	Title             string
	AlternativeTitles []string
	Author            string
	Description       string
	Tags              []string
	Preview           *multipart.FileHeader
	AgeRating         int
	ReleaseYear       int
}

type UpdateMangaRawDTO struct {
	ID                string                `form:"-" bson:"_id" swaggerignore:"true"`
	Title             *string               `form:"title"`
	Author            *string               `form:"author"`
	AlternativeTitles *string               `form:"alternativeTitles" example:"title1,title2,title3 (separate by comma)"`
	Description       *string               `form:"description"`
	Tags              *string               `form:"tags" example:"tag1,tag2,tag3 (separate by comma)"`
	Preview           *multipart.FileHeader `form:"preview" swaggerignore:"true"`
	AgeRating         *int                  `form:"ageRating"`
	ReleaseYear       *int                  `form:"releaseYear"`
	IsPublished       *bool                 `form:"isPublished"`
	Slug              string                `form:"-" swaggerignore:"true"`
}

func (m *UpdateMangaRawDTO) Validate() UpdateMangaDTO {
	var mangaDTO = UpdateMangaDTO{
		ID:                m.ID,
		Title:             m.Title,
		Author:            m.Author,
		AlternativeTitles: nil,
		Description:       m.Description,
		Tags:              nil,
		Preview:           m.Preview,
		AgeRating:         m.AgeRating,
		ReleaseYear:       m.ReleaseYear,
		IsPublished:       m.IsPublished,
		Slug:              m.Slug,
	}

	if m.Tags != nil {
		tags := parseStringToArray(*m.Tags)
		mangaDTO.Tags = &tags
	}

	if m.AlternativeTitles != nil {
		altTitles := parseStringToArray(*m.AlternativeTitles)
		mangaDTO.AlternativeTitles = &altTitles
	}

	return mangaDTO
}

type UpdateMangaDTO struct {
	ID                string
	Title             *string
	Author            *string
	AlternativeTitles *[]string
	Description       *string
	Tags              *[]string
	Preview           *multipart.FileHeader
	PreviewUrl        *string
	AgeRating         *int
	ReleaseYear       *int
	IsPublished       *bool
	Slug              string
}

type RawMangaFilter struct {
	Offset      *int    `form:"offset,default=0"`
	OrderBy     *string `form:"orderBy,default=date"`
	Tags        *string `form:"tags"`
	IsPublished *bool   `form:"isPublished,default=true"`
}

func (f *RawMangaFilter) Validate() (MangaFilter, error) {
	var filter = MangaFilter{
		Offset:      *f.Offset,
		OrderBy:     *f.OrderBy,
		IsPublished: *f.IsPublished,
	}

	if f.Tags != nil {
		tags := parseStringToArray(*f.Tags)
		filter.Tags = tags
	}
	if f.OrderBy != nil {
		if !strings.Contains(orderByValues, *f.OrderBy) {
			return filter, fmt.Errorf("%v: orderBy can be: %v", ErrInvalidFilter.Error(), orderByValues)
		}
	}

	return filter, nil
}

type MangaFilter struct {
	Offset      int
	OrderBy     string
	Tags        []string
	IsPublished bool
}

func parseStringToArray(rawString string) []string {
	var array = make([]string, 0)
	for _, tag := range strings.Split(rawString, ",") {
		if len(tag) > 0 {
			array = append(array, tag)
		}
	}

	return array
}
