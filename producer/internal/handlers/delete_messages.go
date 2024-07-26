package handlers

import (
	"github.com/labstack/echo/v4"
	"net/http"
)

func (h *Handlers) deleteMessages(c echo.Context) error {
	err := h.services.DeleteMessages()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, ErrorResponse{Message: "error deleting messages"})
	}

	return c.NoContent(http.StatusNoContent)
}
