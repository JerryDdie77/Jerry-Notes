package httpserver

import (
	"errors"
	"jerry-notes/internal/service"
	"net/http"

	"github.com/gin-gonic/gin"
)

func (h *Handler) StartRegistration(c *gin.Context) {

	var req struct {
		Name     string `json:"name" binding:"required"`
		Email    string `json:"email" binding:"required"`
		Password string `json:"password" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "bad request",
			"data":  err.Error(),
		})
		return
	}

	ctx := c.Request.Context()

	err := h.app.AuthService.StartRegistration(ctx, req.Name, req.Password, req.Email)
	if err != nil {
		if errors.Is(err, service.ErrInternal) {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		if errors.Is(err, service.ErrEmailTaken) {
			c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
			return
		}

		if errors.Is(err, service.ErrPasswordTooLong) {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
	}

	c.JSON(http.StatusOK, gin.H{"message": "code sent to email"})

}

func (h *Handler) VerifyCode(c *gin.Context) {

	var req struct {
		Email string `json:"email" binding:"required"`
		Code  string `json:"code" binding:"required"`
	}

	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "bad request",
			"data":  err.Error(),
		})
		return
	}

	ctx := c.Request.Context()

	jwtToken, err := h.app.AuthService.VerifyCode(ctx, req.Email, req.Code)

	if err != nil {

		if errors.Is(err, service.ErrInternal) {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
			return
		}

		if errors.Is(err, service.ErrNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}

		if errors.Is(err, service.ErrCodeExpired) {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}

		if errors.Is(err, service.ErrInvalidCode) {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
	}

	c.SetCookie("auth_token", jwtToken, 7*24*60*60, "/", "", true, true)

	c.JSON(http.StatusCreated, gin.H{"message": "registered and logged in"})

}
