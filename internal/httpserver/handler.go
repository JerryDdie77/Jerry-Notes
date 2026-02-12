package httpserver

import (
	"jerry-notes/internal/app"
	"net/http"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	app *app.App
}

func NewHandler(app *app.App) *Handler {
	return &Handler{
		app: app,
	}
}

func (h *Handler) Ping(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}
