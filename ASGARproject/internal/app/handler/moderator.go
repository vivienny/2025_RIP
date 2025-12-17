package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// GetAllUsers возвращает всех пользователей (для модератора)
func (h *Handler) GetAllUsers(ctx *gin.Context) {
	// Временная заглушка
	ctx.JSON(http.StatusOK, gin.H{
		"message": "Метод GetAllUsers будет реализован",
		"users":   []string{},
		"count":   0,
	})
}

// UpdateUserRole обновляет роль пользователя
func (h *Handler) UpdateUserRole(ctx *gin.Context) {
	idStr := ctx.Param("id")

	ctx.JSON(http.StatusOK, gin.H{
		"message": "Метод UpdateUserRole будет реализован",
		"user_id": idStr,
		"role":    "not implemented",
	})
}
