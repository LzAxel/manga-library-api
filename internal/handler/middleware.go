package handler

import (
	"errors"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

const (
	authHeader = "Authorization"
)

func (h *Handler) userIdentity(ctx *gin.Context) {
	h.logger.Debugln("user identification")

	header := ctx.GetHeader("Authorization")
	if header == "" {
		ErrorResponse(ctx, http.StatusUnauthorized, "empty auth header")
		return
	}

	headerParts := strings.Split(header, " ")
	if len(headerParts) != 2 || headerParts[0] != "Bearer" {
		ErrorResponse(ctx, http.StatusUnauthorized, "invalid auth header")
		return
	}

	headerToken := headerParts[1]
	token, err := h.services.JWTMangaer.Parse(headerToken)
	if err != nil {
		ErrorResponse(ctx, http.StatusUnauthorized, "invalid token")
		return
	}

	claims, err := h.services.JWTMangaer.Claims(token)
	if err != nil {
		ErrorResponse(ctx, http.StatusUnauthorized, "invalid token claims")
		return
	}

	ctx.Set("userId", claims["sub"].(string))
}

func (h *Handler) getUserId(ctx *gin.Context) (string, error) {
	userId, ok := ctx.Value("userId").(string)
	if userId == "" || !ok {
		return "", errors.New("failed to get user id")
	}

	return userId, nil
}
