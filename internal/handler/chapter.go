package handler

import (
	"manga-library/internal/domain"
	"net/http"
	"path/filepath"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
)

const zipExtansion = ".zip"

func (h *Handler) uploadChapter(ctx *gin.Context) {
	h.logger.Debugln("uploading chapter")
	var chapterDTO domain.UploadChapterDTO

	uploaderID, err := h.getUserId(ctx)
	if err != nil {
		ErrorResponse(ctx, http.StatusBadRequest, err.Error())
		return
	}
	roles, err := h.getUserRoles(ctx, uploaderID)
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

	if err := h.services.Manga.UploadChapter(ctx, chapterDTO, roles); err != nil {
		ErrorResponse(ctx, http.StatusInternalServerError, err.Error())
		return
	}
}

func (h *Handler) deleteChapter(ctx *gin.Context) {
	h.logger.Debugln("deleting chapter")
	uploaderID, err := h.getUserId(ctx)
	if err != nil {
		ErrorResponse(ctx, http.StatusBadRequest, err.Error())
		return
	}
	roles, err := h.getUserRoles(ctx, uploaderID)
	if err != nil {
		ErrorResponse(ctx, http.StatusBadRequest, err.Error())
		return
	}

	slug := ctx.Param("mangaSlug")
	volume, err := strconv.Atoi(ctx.Param("volume"))
	if err != nil {
		ErrorResponse(ctx, http.StatusBadRequest, "invalid input")
		return
	}
	number, err := strconv.Atoi(ctx.Param("number"))
	if err != nil {
		ErrorResponse(ctx, http.StatusBadRequest, "invalid input")
		return
	}

	if slug == "" || volume == 0 || number == 0 {
		ErrorResponse(ctx, http.StatusBadRequest, "invalid input")
		return
	}

	err = h.services.Manga.DeleteChapter(ctx, domain.DeleteChapterDTO{
		MangaSlug:  slug,
		UploaderID: uploaderID,
		Volume:     volume,
		Number:     number,
	}, roles)
	if err != nil {
		ErrorResponse(ctx, http.StatusBadRequest, err.Error())
		return
	}

	ctx.Writer.WriteHeader(http.StatusOK)
}
