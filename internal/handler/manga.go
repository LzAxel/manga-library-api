package handler

import (
	"errors"
	"manga-library/internal/domain"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/mongo"
)

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
		if errors.Is(err, mongo.ErrNoDocuments) {
			ErrorResponse(c, http.StatusNotFound, "manga not found")
			return
		}
		h.logger.Errorln(err)
		ErrorResponse(c, http.StatusInternalServerError, "failed to getting manga")
		return
	}

	c.JSON(http.StatusOK, manga)
}

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
}
