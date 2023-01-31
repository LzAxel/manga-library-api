package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// Heartbeat
// @Summary Heartbeat metric
// @Tags Metrics
// @Success 204
// @Failure 400
// @Router /api/heartbeat [get]
func (h *Handler) heartbeat(c *gin.Context) {
	c.JSON(http.StatusOK, map[string]interface{}{"status": "ok"})
}
