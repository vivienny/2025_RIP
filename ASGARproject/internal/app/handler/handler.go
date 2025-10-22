package handler

import (
	"decode/internal/app/repository"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

type Handler struct {
	Repository *repository.Repository
}

func NewHandler(r *repository.Repository) *Handler {
	return &Handler{
		Repository: r,
	}
}

// RegisterHandler Функция, в которой мы отдельно регистрируем маршруты
func (h *Handler) RegisterHandler(router *gin.Engine) {
	router.GET("/", h.GetAsgarServices)
	router.GET("/service/:id", h.GetAsgarService)
	router.GET("/miniplane", h.GetMiniplane)
	router.POST("/remove-from-miniplane", h.RemoveFromMiniplane)
	router.POST("/clear-miniplane", h.ClearMiniplane) // ← ДОБАВЛЯЕМ ОЧИСТКУ
}

// RegisterStatic То же самое, что и с маршрутами, регистрируем статику
func (h *Handler) RegisterStatic(router *gin.Engine) {
	router.LoadHTMLGlob("templates/*")
	router.Static("/static", "./resources")
}

func (h *Handler) errorHandler(ctx *gin.Context, errorStatusCode int, err error) {
	logrus.Error(err.Error())
	ctx.JSON(errorStatusCode, gin.H{
		"status":      "error",
		"description": err.Error(),
	})
}
