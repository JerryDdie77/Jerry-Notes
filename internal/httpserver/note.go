package httpserver

import (
	"errors"
	"jerry-notes/internal/service"
	"net/http"

	"github.com/gin-gonic/gin"
)

func (h *Handler) CreateNote(c *gin.Context) {

	var req struct {
		Title   string `json:"title" binding:"required"`
		Content string `json:"content" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "bad request",
			"data":  err.Error(),
		})
		return
	}

	userID, ok := getUserID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "not authorized"})
		return
	}

	ctx := c.Request.Context()

	noteID, err := h.app.NoteService.CreateNote(ctx, userID, req.Title, req.Content)

	if err != nil {

		if errors.Is(err, service.ErrInternal) {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		if errors.Is(err, service.ErrEmptyTitle) {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "successful creating the note",
		"note_id": noteID,
	})

}
