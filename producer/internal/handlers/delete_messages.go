package handlers

import (
	"Messaggio/producer/internal/model"
	"errors"
	"github.com/labstack/echo/v4"
	"net/http"
)

// deleteMessages
// @Summary Delete messages
// @Tags messages
// @Description Delete all messages
// @ID delete-messages
// @Success 204
// @Failure 500 {object} ErrorResponse
// @Router /messages [delete]
func (h *Handlers) deleteMessages(c echo.Context) error {
	err := h.services.DeleteMessages()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, ErrorResponse{Message: "error deleting messages"})
	}

	return c.NoContent(http.StatusNoContent)
}

// deleteMessage
// @Summary Delete message
// @Tags messages
// @Description Delete message by id
// @ID delete-message
// @Param id path string true "Message ID"
// @Success 204
// @Failure 400 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
func (h *Handlers) deleteMessage(c echo.Context) error {
	id := c.Param("id")
	if id == "" {
		return c.JSON(http.StatusBadRequest, &ErrorResponse{Message: "id is required"})
	}

	err := h.services.DeleteMessage(id)
	if err != nil {
		if errors.Is(err, model.ErrNotFound) {
			return c.JSON(http.StatusNotFound, &ErrorResponse{Message: "message not found"})
		}
		return c.JSON(http.StatusInternalServerError, &ErrorResponse{Message: "error deleting message"})
	}

	return c.NoContent(http.StatusNoContent)
}
