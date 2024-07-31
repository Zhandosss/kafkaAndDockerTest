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

// getMessage
// @Summary Get message
// @Tags messages
// @Description Get message by id
// @ID get-message
// @Accept json
// @Produce json
// @Param id path string true "Message ID"
// @Success 200 {object} GetMessageResponse
// @Failure 400 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /messages/{id} [get]
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

// getMessages
// @Summary Get messages
// @Tags messages
// @Description Get all messages
// @ID get-messages
// @Accept json
// @Produce json
// @Success 200 {object} GetMessagesResponse
// @Failure 500 {object} ErrorResponse
// @Router /messages [get]
func (h *Handlers) getMessages(c echo.Context) error {
	messages, err := h.services.GetMessages()
	if err != nil {
		log.Error().Err(err).Msg("failed to get messages")
		return c.JSON(http.StatusInternalServerError, &ErrorResponse{Message: "failed to get messages"})
	}

	//TODO make struct to return
	return c.JSON(http.StatusOK, GetMessagesResponse{Messages: messages})
}
