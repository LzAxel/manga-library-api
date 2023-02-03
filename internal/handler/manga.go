package handler

import (
	"errors"
	"manga-library/internal/domain"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

// Create Manga
// @Summary Create Manga
// @Tags Manga
// @Accept json
// @Security BearerAuth
// @Param manga body domain.CreateMangaDTO true "Add manga"
// @Success 200 {object} string "Created manga ID"
// @Failure 400 "Invalid input"
// @Failure 500
// @Router /api/manga [post]
func (h *Handler) createManga(c *gin.Context) {
	var mangaDTO domain.CreateMangaDTO

	userId, err := h.getUserId(c)
	if err != nil {
		h.logger.Errorln(err)
		ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}
	h.logger.WithFields(logrus.Fields{"userId": userId}).Debugln("creating manga")

	if err := c.BindJSON(&mangaDTO); err != nil {
		h.logger.Errorln(err)
		ErrorResponse(c, http.StatusBadRequest, "invalid input")
		return
	}

	id, err := h.services.Manga.Create(c, userId, mangaDTO)
	if err != nil {
		h.logger.Errorln(err)
		ErrorResponse(c, http.StatusInternalServerError, "failed to save manga")
		return
	}

	c.JSON(http.StatusCreated, map[string]string{"id": id})
}

// Get Latest Manga
// @Summary Get Latest Manga
// @Tags Manga
// @Success 200 {array} domain.Manga
// @Failure 500
// @Router /api/manga/latest [get]
func (h *Handler) getLatestManga(c *gin.Context) {
	// TODO: context with timeout for all service layers
	h.logger.Debugln("getting latest manga")

	result, err := h.services.Manga.GetLatest(c)
	if err != nil {
		h.logger.Errorln(err)
		ErrorResponse(c, http.StatusInternalServerError, "failed to get latest manga")
		return
	}

	c.JSON(http.StatusOK, result)
}

// Get Manga
// @Summary Get Manga by ID or Slug
// @Tags Manga
// @Success 200 {object} domain.Manga
// @Failure 400 "Invalid input (One of params is required)"
// @Failure 404 "Manga not found"
// @Failure 500
// @Param id query string false "Manga ID"
// @Param slug query string false "Manga slug"
// @Router /api/manga [get]
func (h *Handler) getManga(c *gin.Context) {
	var manga domain.Manga

	slug := c.Query("slug")
	id := c.Query("id")

	h.logger.WithFields(logrus.Fields{"id": id, "slug": slug}).Debugln("getting manga")

	if slug == "" && id == "" {
		ErrorResponse(c, http.StatusBadRequest, "slug and id params are empty")
		return
	}
	if slug != "" && id != "" {
		ErrorResponse(c, http.StatusBadRequest, "allowed only one of params: id, slug")
		return
	}

	getMangaDTO := domain.GetMangaDTO{
		Id:   id,
		Slug: slug,
	}

	manga, err := h.services.Manga.Get(c, getMangaDTO)

	if err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			ErrorResponse(c, http.StatusNotFound, "manga not found")
			return
		}
		h.logger.Errorln(err)
		ErrorResponse(c, http.StatusInternalServerError, "failed to getting manga")
		return
	}

	c.JSON(http.StatusOK, manga)
}

// Delete Manga
// @Summary Delete Manga by ID
// @Security BearerAuth
// @Tags Manga
// @Success 204
// @Failure 400
// @Failure 500
// @Param id path string true "Manga ID"
// @Router /api/manga/{id} [delete]
func (h *Handler) deleteManga(c *gin.Context) {
	var mangaId = c.Param("id")

	userId, err := h.getUserId(c)
	if err != nil {
		h.logger.Errorln(err)
		ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	h.logger.WithFields(logrus.Fields{"mangaId": mangaId}).Debugln("deleting manga")

	if err := h.services.Manga.Delete(c, userId, mangaId); err != nil {
		h.logger.Errorln(err)
		ErrorResponse(c, http.StatusInternalServerError, "failed to delete manga")
		return
	}

	c.JSON(http.StatusNoContent, nil)
}

// Update Manga
// @Summary Update Manga by ID
// @Security BearerAuth
// @Accept json
// @Tags Manga
// @Success 200
// @Failure 400
// @Failure 500
// @Param manga body domain.UpdateMangaDTO true "Update manga"
// @Param id path string true "Manga ID"
// @Router /api/manga/{id} [patch]
func (h *Handler) updateManga(c *gin.Context) {
	h.logger.Debugln("updating manga")
	var mangaId = c.Param("id")

	userId, err := h.getUserId(c)
	if err != nil {
		h.logger.Errorln(err)
		ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	var mangaDTO domain.UpdateMangaDTO

	if err := c.BindJSON(&mangaDTO); err != nil {
		h.logger.Errorln(err)
		ErrorResponse(c, http.StatusBadRequest, "invalid input")
		return
	}
	mangaDTO.Id = mangaId
	err = h.services.Manga.Update(c, userId, mangaDTO)
	if err != nil {
		h.logger.Errorln(err)
		ErrorResponse(c, http.StatusInternalServerError, "failed to update manga")
		return
	}

	c.Writer.WriteHeader(http.StatusOK)
}
