package domain

import (
	"errors"
	"time"
)

var (
	ErrNotFound    = errors.New("not found")
	ErrNotTheOwner = errors.New("you are not the owner")
	ErrNotEditor   = errors.New("you are not editor")

	ErrWrongAuthCreditionals = errors.New("wrong password or username")
	ErrUsernameExists        = errors.New("username already exists")
)

type User struct {
	ID           string    `json:"_id" bson:"_id"`
	Username     string    `json:"username" bson:"username"`
	PasswordHash []byte    `json:"-" bson:"passwordHash"`
	IsEditor     bool      `json:"isEditor" bson:"isEditor"`
	IsAdmin      bool      `json:"isAdmin" bson:"isAdmin"`
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
