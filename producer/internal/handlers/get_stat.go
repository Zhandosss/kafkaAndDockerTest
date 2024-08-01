package handlers

import (
	"Messaggio/producer/internal/model"
	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog/log"
	"net/http"
)

type GetStatResponse struct {
	Stats map[string]*model.ByDays `json:"stats"`
}

// getStatsByDays
// @Summary Get stats
// @Tags stats
// @Description Get stats by days
// @ID get-stats
// @Produce json
// @Success 200 {object} GetStatResponse
// @Failure 500 {object} ErrorResponse
// @Router /statistic/days [get]
func (h *Handlers) getStatsByDays(c echo.Context) error {
	requestID := c.Response().Header().Get(echo.HeaderXRequestID)

	log.Info().Msgf("RequestID: %s, Get stats request", requestID)

	stat, err := h.services.GetStatsByDays()
	if err != nil {
		log.Error().Err(err).Msgf("RequestID: %s, Failed to get stats", requestID)
		return c.JSON(http.StatusInternalServerError, &ErrorResponse{Message: "failed to get stats"})
	}

	log.Info().Msgf("RequestID: %s, Stats: %+v", requestID, stat)
	return c.JSON(http.StatusOK, &GetStatResponse{Stats: stat})
}
