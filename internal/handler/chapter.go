package handler

import (
	"manga-library/internal/domain"
	"net/http"
	"path/filepath"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
)

const zipExtansion = ".zip"

func (h *Handler) uploadChapter(ctx *gin.Context) {
	var chapterDTO domain.UploadChapterDTO

	uploaderID, err := h.getUserId(ctx)
	if err != nil {
		ErrorResponse(ctx, http.StatusBadRequest, err.Error())
		return
	}

	// TODO: unified error
	if err := ctx.ShouldBindWith(&chapterDTO, binding.FormMultipart); err != nil {
		ErrorResponse(ctx, http.StatusBadRequest, "invalid input")
		return
	}
	chapterDTO.UploaderID = uploaderID

	if filepath.Ext(chapterDTO.File.Filename) != zipExtansion {
		ErrorResponse(ctx, http.StatusBadRequest, "invalid file extansion (should be zip)")
		return
	}

	h.logger.Debugln("binded successfully")
	if err := h.services.Manga.UploadChapter(ctx, chapterDTO); err != nil {
		ErrorResponse(ctx, http.StatusInternalServerError, err.Error())
		return
	}

	h.logger.Debugln(chapterDTO)
}
