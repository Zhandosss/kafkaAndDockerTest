package handlers

import (
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
