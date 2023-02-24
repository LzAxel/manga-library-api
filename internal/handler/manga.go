package handler

import (
	"errors"
	"manga-library/internal/domain"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/sirupsen/logrus"
)

// Create Manga
// @Summary Create Manga
// @Tags Manga
// @Accept mpfd
// @Security BearerAuth
// @Param file formData domain.CreateMangaRawDTO true "request"
// @Param file formData file true "Manga preview image (jpg, jpeg)"
// @Success 200 {object} string "Created manga ID"
// @Failure 400 "Invalid input"
// @Failure 500
// @Router /api/manga/ [post]
func (h *Handler) createManga(c *gin.Context) {
	var mangaRawDTO domain.CreateMangaRawDTO
	var mangaDTO domain.CreateMangaDTO

	userId, err := h.getUserId(c)
	if err != nil {
		h.logger.Errorln(err)
		ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}
	h.logger.WithFields(logrus.Fields{"userId": userId}).Debugln("creating manga")

	if err := c.ShouldBindWith(&mangaRawDTO, binding.FormMultipart); err != nil {
		h.logger.Errorln(err)
		ErrorResponse(c, http.StatusBadRequest, "invalid input")
		return
	}
	mangaDTO = mangaRawDTO.Validate()

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

	if len(result) == 0 {
		ErrorResponse(c, http.StatusNotFound, domain.ErrNotFound.Error())
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
	var err error

	slug := c.Query("slug")
	id := c.Query("id")

	h.logger.WithFields(logrus.Fields{"id": id, "slug": slug}).Debugln("getting manga")

	if (slug != "" && id != "") || (slug == "" && id == "") {
		ErrorResponse(c, http.StatusBadRequest, "invalid input")
		return
	}

	if slug != "" {
		manga, err = h.services.Manga.GetBySlug(c, slug)
	} else {
		manga, err = h.services.Manga.GetByID(c, id)
	}

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

// Get Manga By Tags
// @Summary Get Manga by tags
// @Tags Manga
// @Success 200 {object} []domain.Manga
// @Failure 400 "Invalid input (Invalid offset)"
// @Failure 404 "Manga not found"
// @Failure 500
// @Param tags query string true "Tags"
// @Param offset query string false "Offset"
// @Router /api/manga/filter [get]
func (h *Handler) getMangaByFilter(c *gin.Context) {
	var (
		rawFilter domain.RawMangaFilter
		filter    domain.MangaFilter
		err       error
	)

	h.logger.Debugln("getting manga by filter")

	if err = c.ShouldBindQuery(&rawFilter); err != nil {
		ErrorResponse(c, http.StatusBadRequest, domain.ErrInvalidFilter.Error())
		return
	}
	filter, err = rawFilter.Validate()
	if err != nil {
		ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	h.logger.Debugf("%v %v %v %v", filter.Offset, filter.OrderBy,
		filter.Tags, filter.IsPublished)

	// TODO: make something with kostb|l'
	if !filter.IsPublished {
		h.userIdentity(c)
		userId, err := h.getUserId(c)
		if err == nil {
			roles, err := h.getUserRoles(c, userId)
			if err != nil {
				return
			}
			if roles.IsAdmin && !filter.IsPublished {
				filter.IsPublished = false
			} else {
				filter.IsPublished = true
			}

		} else {
			return
		}
	}

	manga, err := h.services.Manga.GetByFilter(c, filter)
	if err != nil {
		h.logger.Errorln(err)
		ErrorResponse(c, http.StatusInternalServerError, "failed to getting manga by filter")
		return
	}
	if len(manga) == 0 {
		ErrorResponse(c, http.StatusNotFound, domain.ErrNotFound.Error())
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
// @Accept mpfd
// @Tags Manga
// @Success 200
// @Failure 400
// @Failure 500
// @Param file formData domain.UpdateMangaRawDTO true "request"
// @Param preview formData file false "Manga preview image (jpg, jpeg)"
// @Param id path string true "Manga ID"
// @Router /api/manga/{id} [patch]
func (h *Handler) updateManga(c *gin.Context) {
	h.logger.Debugln("updating manga")
	var mangaId = c.Param("id")
	var mangaRawDTO domain.UpdateMangaRawDTO
	var mangaDTO domain.UpdateMangaDTO

	userId, err := h.getUserId(c)
	if err != nil {
		h.logger.Errorln(err)
		ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	roles, err := h.getUserRoles(c, userId)
	if err != nil {
		h.logger.Errorln(err)
		ErrorResponse(c, http.StatusInternalServerError, "failed to get user roles")
		return
	}

	if err := c.ShouldBindWith(&mangaRawDTO, binding.FormMultipart); err != nil {
		h.logger.Errorln(err)
		ErrorResponse(c, http.StatusBadRequest, "invalid input")
		return
	}
	mangaDTO = mangaRawDTO.Validate()

	mangaDTO.ID = mangaId
	err = h.services.Manga.Update(c, userId, roles, mangaDTO)
	if err != nil {
		h.logger.Errorln(err)
		ErrorResponse(c, http.StatusInternalServerError, "failed to update manga")
		return
	}

	c.Writer.WriteHeader(http.StatusOK)
}
