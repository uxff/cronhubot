package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/uxff/cronhubot/pkg/checker"
	"github.com/uxff/cronhubot/pkg/render"
)

type HealthzHandler struct {
	checkers map[string]checker.Checker
}

func NewHealthzHandler(checkers map[string]checker.Checker) *HealthzHandler {
	return &HealthzHandler{checkers}
}

func (h *HealthzHandler) HealthzIndex(c *gin.Context) {
	payload := make(map[string]bool)

	for k, v := range h.checkers {
		payload[k] = v.IsAlive()
	}

	render.JSON(c.Writer, http.StatusOK, payload)
}
