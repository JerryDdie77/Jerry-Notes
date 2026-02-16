package httpserver

import (
	"errors"
	"jerry-notes/internal/service"
	"net/http"
	"strconv"

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

func (h *Handler) GetNote(c *gin.Context) {

	userID, ok := getUserID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "not authorized"})
		return
	}

	noteIDStr := c.Param("id")

	noteID, err := strconv.Atoi(noteIDStr)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "not found"})
		return
	}

	ctx := c.Request.Context()

	note, err := h.app.NoteService.GetNote(ctx, int64(noteID), userID)
	if err != nil {
		switch {
		case errors.Is(err, service.ErrNotFound):
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		case errors.Is(err, service.ErrInternal):
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		case errors.Is(err, service.ErrForbidden):
			c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
			return
		}
	}

	c.JSON(http.StatusOK, note)
}

func (h *Handler) DeleteNote(c *gin.Context) {
	userID, ok := getUserID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "not authorized"})
		return
	}

	noteIDStr := c.Param("id")

	noteID, err := strconv.Atoi(noteIDStr)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "not found"})
		return
	}

	ctx := c.Request.Context()

	err = h.app.NoteService.DeleteNote(ctx, userID, int64(noteID))

	if err != nil {
		switch {
		case errors.Is(err, service.ErrNotFound):
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		case errors.Is(err, service.ErrInternal):
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		case errors.Is(err, service.ErrForbidden):
			c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
			return
		}
	}

	c.JSON(http.StatusOK, gin.H{"message": "successfully deleted"})

}

func (h *Handler) ListNotes(c *gin.Context) {

	userID, ok := getUserID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "not authorized"})
		return
	}

	ctx := c.Request.Context()

	notes, err := h.app.NoteService.ListNotes(ctx, userID)
	if errors.Is(err, service.ErrInternal) {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, notes)
}
