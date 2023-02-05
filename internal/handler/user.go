package handler

import (
	"manga-library/internal/domain"
	"net/http"

	"github.com/gin-gonic/gin"
)

func (h *Handler) getUserByUsername(ctx *gin.Context) {
	username := ctx.Param("username")
	if username == "" {
		ErrorResponse(ctx, http.StatusBadRequest, "username is empty")
		return
	}
	h.logger.Debugln(username)

	user, err := h.services.User.GetByUsername(ctx, username)
	if err != nil {
		switch err {
		case domain.ErrNotFound:
			ErrorResponse(ctx, http.StatusBadRequest, err.Error())
			return
		default:
			ErrorResponse(ctx, http.StatusInternalServerError, "failed to get user")
			return
		}
	}

	ctx.JSON(http.StatusOK, user)
}
func (h *Handler) updateUser(ctx *gin.Context) {

}
func (h *Handler) deleteUser(ctx *gin.Context) {

}
