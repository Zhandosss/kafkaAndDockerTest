package handlers

import (
	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog/log"
	"net/http"
)

func (h *Handlers) getMessage(c echo.Context) error {
	id := c.Param("id")
	if id == "" {
		return c.JSON(http.StatusBadRequest, &ErrorResponse{Message: "id is required"})
	}

	message, err := h.services.GetMessage(id)
	if err != nil {
		log.Error().Err(err).Msg("failed to get message")
		return c.JSON(http.StatusInternalServerError, &ErrorResponse{Message: "failed to get message"})
	}

	//TODO make struct to return
	return c.JSON(http.StatusOK, message)
}

func (h *Handlers) getMessages(c echo.Context) error {
	messages, err := h.services.GetMessages()
	if err != nil {
		log.Error().Err(err).Msg("failed to get messages")
		return c.JSON(http.StatusInternalServerError, &ErrorResponse{Message: "failed to get messages"})
	}

	//TODO make struct to return
	return c.JSON(http.StatusOK, messages)
}
