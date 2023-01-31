package handler

import (
	"context"
	"errors"
	"manga-library/internal/domain"
	appErrors "manga-library/pkg/errors"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

const (
	responseTimeout = time.Second * 30
)

func (h *Handler) signUp(c *gin.Context) {
	h.logger.Debugln("signing up user")
	var userDTO domain.CreateUserDTO

	if err := c.BindJSON(&userDTO); err != nil {
		return
	}

	if err := h.services.Authorization.SignUp(c, userDTO); err != nil {
		h.logger.Errorln(err)
		if errors.Is(err, appErrors.ErrUsernameExists) {
			ErrorResponse(c, http.StatusConflict, err.Error())
		} else {
			ErrorResponse(c, http.StatusInternalServerError, err.Error())
		}
		return
	}

	c.Writer.WriteHeader(http.StatusCreated)
}
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
		ErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	h.logger.WithFields(logrus.Fields{"token": token}).Debugf("user authorized successfully")

	c.JSON(200, map[string]string{"token": token})
}
