package handlers

import (
	"Messaggio/producer/internal/model"
	"errors"
	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog/log"
	"net/http"
)

type GetMessageResponse struct {
	Message *model.Message `json:"message"`
}

func (h *Handlers) getMessage(c echo.Context) error {
	id := c.Param("id")
	if id == "" {
		return c.JSON(http.StatusBadRequest, &ErrorResponse{Message: "id is required"})
	}

	message, err := h.services.GetMessage(id)
	if err != nil {
		if errors.Is(err, model.ErrNotFound) {
			return c.JSON(http.StatusNotFound, &ErrorResponse{Message: "message not found"})
		}
		log.Error().Err(err).Msg("failed to get message")
		return c.JSON(http.StatusInternalServerError, &ErrorResponse{Message: "failed to get message"})
	}

	//TODO make struct to return
	return c.JSON(http.StatusOK, GetMessageResponse{Message: message})
}

type GetMessagesResponse struct {
	Messages []*model.Message `json:"messages"`
}

func (h *Handlers) getMessages(c echo.Context) error {
	messages, err := h.services.GetMessages()
	if err != nil {
		log.Error().Err(err).Msg("failed to get messages")
		return c.JSON(http.StatusInternalServerError, &ErrorResponse{Message: "failed to get messages"})
	}

	//TODO make struct to return
	return c.JSON(http.StatusOK, GetMessagesResponse{Messages: messages})
}
