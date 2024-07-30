package handlers

import (
	"Messaggio/producer/internal/model"
	"github.com/labstack/echo/v4"
	"net/http"
)

type GetStatResponse struct {
	Stats map[string]*model.ByDays `json:"stats"`
}

func (h *Handlers) getStatsByDays(c echo.Context) error {
	stat, err := h.services.GetStatsByDays()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, &ErrorResponse{Message: "failed to get stats"})
	}

	return c.JSON(http.StatusOK, &GetStatResponse{Stats: stat})
}
