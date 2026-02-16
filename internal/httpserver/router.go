package httpserver

import "github.com/gin-gonic/gin"

func NewRouter(h *Handler) *gin.Engine {
	r := gin.Default()
	r.GET("/health", h.Ping)

	api := r.Group("/api")

	authGroup := api.Group("/auth")
	{
		authGroup.POST("/register", h.StartRegistration)
		authGroup.POST("/verify-code", h.VerifyCode)
		authGroup.POST("/login", h.Login)
	}

	protected := api.Group("/protected")

	protected.Use(h.JWTMiddleware())

	protected.POST("/notes", h.CreateNote)
	protected.GET("/notes/:id", h.GetNote)

	return r
}
