package domain

import (
	"errors"
	"time"
)

var (
	ErrNotFound   = errors.New("not found")
	ErrNotAnOwner = errors.New("not an owner of manga")

	ErrWrongAuthCreditionals = errors.New("wrong password or username")
	ErrUsernameExists        = errors.New("username already exists")
)

type User struct {
	ID           string    `json:"_id" bson:"_id"`
	Username     string    `json:"username" bson:"username"`
	PasswordHash string    `json:"passwordHash" bson:"passwordHash"`
	IsEditor     bool      `json:"isEditor" bson:"isEditor"`
	CreatedAt    time.Time `json:"createdAt" bson:"createdAt"`
}

type AdminUser struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type CreateUserDTO struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type LoginUserDTO struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}
