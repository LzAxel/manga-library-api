package domain

import (
	"errors"
	"time"
)

var (
	ErrNotFound    = errors.New("not found")
	ErrFailedToGet = errors.New("failed to get")
	ErrNotTheOwner = errors.New("not the owner of note")
)

type User struct {
	Id           string    `json:"_id" bson:"_id"`
	Username     string    `json:"username" bson:"username"`
	PasswordHash string    `json:"passwordHash" bson:"passwordHash"`
	IsEditor     bool      `json:"isEditor" bson:"isEditor"`
	CreatedAt    time.Time `json:"createdAt" bson:"createdAt"`
}

type CreateUserDTO struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type LoginUserDTO struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}
