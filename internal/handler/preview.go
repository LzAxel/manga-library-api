package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

const previewUploadDir = "C:/Users/lzaxel/code/go_projects/manga-library/upload/"

func (h *Handler) uploadMangaPreview(c *gin.Context) {
	file, _ := c.FormFile("file")

	mp, err := file.Open()
	if err != nil {
		h.logger.Errorln(err)
		ErrorResponse(c, http.StatusInternalServerError, "failed to open preview")
		return
	}
	defer mp.Close()

	previewUrl, err := h.services.Preview.Create(c, mp, file.Filename, "fasdfsdfsdf")
	if err != nil {
		h.logger.Errorln(err)
		ErrorResponse(c, http.StatusInternalServerError, "failed to save preview")
		return
	}

	c.JSON(http.StatusCreated, map[string]string{"url": previewUrl})
}

func (h *Handler) deleteMangaPreview(c *gin.Context) {
	previewId := c.Param("id")

	if err := h.services.Preview.Delete(c, previewId); err != nil {
		ErrorResponse(c, http.StatusInternalServerError, "failed to delete preview")
		return
	}

	c.JSON(http.StatusNoContent, nil)
}
