package httpserver

import "github.com/gin-gonic/gin"

func NewRouter(h *Handler) *gin.Engine {
	r := gin.Default()

	r.GET("/health", h.Ping)

	return r
}
