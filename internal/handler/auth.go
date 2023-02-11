package handler

import (
	"context"
	"errors"
	"manga-library/internal/domain"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

const (
	responseTimeout = time.Second * 30
)

// Sign Up
// @Summary Sign up
// @Description Sign up via username and password
// @Tags Authorization
// @Accept  json
// @Param auth body domain.CreateUserDTO true "Auth Sign Up Input"
// @Success 201
// @Failure 409 "Username is already exists"
// @Failure 500
// @Router /auth/sign-up [post]
func (h *Handler) signUp(c *gin.Context) {
	h.logger.Debugln("signing up user")
	var userDTO domain.CreateUserDTO

	if err := c.BindJSON(&userDTO); err != nil {
		return
	}

	if err := h.services.Authorization.SignUp(c, userDTO); err != nil {
		h.logger.Errorln(err)
		if errors.Is(err, domain.ErrUsernameExists) {
			ErrorResponse(c, http.StatusConflict, err.Error())
		} else {
			ErrorResponse(c, http.StatusInternalServerError, err.Error())
		}
		return
	}

	c.Writer.WriteHeader(http.StatusCreated)
}

// Sign In
// @Summary Sign In
// @Description Sign in via username and password
// @Tags Authorization
// @Accept  json
// @Produce  json
// @Param auth body domain.LoginUserDTO true "Auth Sign In Input"
// @Success 200 {object} string "token"
// @Failure 400 "Invalid input"
// @Failure 500
// @Router /auth/sign-in [post]
func (h *Handler) signIn(c *gin.Context) {
	h.logger.Debugln("signing in user")
	ctx, cancel := context.WithTimeout(c, responseTimeout)
	defer cancel()

	userDTO := domain.LoginUserDTO{}

	if err := c.BindJSON(&userDTO); err != nil {
		ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	token, err := h.services.Authorization.SignIn(ctx, userDTO)
	if err != nil {
		ErrorResponse(c, http.StatusNotFound, err.Error())
		return
	}

	h.logger.WithFields(logrus.Fields{"token": token}).Debugf("user authorized successfully")

	c.JSON(200, map[string]string{"token": token})
}
