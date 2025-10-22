package handler

import (
	"decode/internal/app/api" // ← ИМПОРТ API
	"decode/internal/app/repository"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

type Handler struct {
	Repository *repository.Repository
	API        *api.MiniplaneHandler // ← API хендлер для корзины
	Services   *api.ServicesHandler  // ← ДОБАВИЛ СЕРВИСЫ
	AuthAvius  *api.AuthAviusHandler
}

func NewHandler(r *repository.Repository) *Handler {
	return &Handler{
		Repository: r,
		API:        api.NewMiniplaneHandler(r), // ← ИНИЦИАЛИЗАЦИЯ API корзины
		Services:   api.NewServicesHandler(r),  // ← ДОБАВИЛ ИНИЦИАЛИЗАЦИЮ СЕРВИСОВ
		AuthAvius:  api.NewAuthAviusHandler(r), // ←
	}
}

// RegisterHandler Функция, в которой мы отдельно регистрируем маршруты
func (h *Handler) RegisterHandler(router *gin.Engine) {
	// HTML маршруты
	router.GET("/", h.GetAsgarServices)
	router.GET("/service/:id", h.GetAsgarService)
	router.GET("/miniplane", h.GetMiniplane)
	router.POST("/remove-from-miniplane", h.RemoveFromMiniplane)
	router.POST("/clear-miniplane", h.ClearMiniplane)

	// API маршруты (REST)
	apiGroup := router.Group("/api")
	{
		// Miniplane (Корзина) - через API хендлер
		apiGroup.GET("/miniplane/icon", h.API.GetMiniplaneIconAPI)
		apiGroup.GET("/miniplane", h.API.GetMiniplaneAPI)
		apiGroup.GET("/miniplane/count", h.API.GetMiniplaneCountAPI)
		apiGroup.DELETE("/miniplane/items/:id", h.API.RemoveFromMiniplaneAPI)
		apiGroup.PUT("/miniplane/items/:id", h.API.UpdateMiniplaneItemAPI)
		apiGroup.DELETE("/miniplane/clear", h.API.ClearMiniplaneAPI)

		// Services (Услуги) - ← ДОБАВИЛ ЭТОТ БЛОК
		apiGroup.GET("/services", h.Services.GetServices)
		apiGroup.GET("/services/:id", h.Services.GetService)
		apiGroup.POST("/services", h.Services.CreateService)
		apiGroup.PUT("/services/:id", h.Services.UpdateService)
		apiGroup.DELETE("/services/:id", h.Services.DeleteService)
		apiGroup.POST("/services/:id/image", h.Services.UploadServiceImage)

		apiGroup.POST("/auth/register", h.AuthAvius.Register)
		apiGroup.POST("/auth/login", h.AuthAvius.Login)
		apiGroup.POST("/auth/logout", h.AuthAvius.Logout)
		apiGroup.GET("/auth/me", h.AuthAvius.GetMe)
		apiGroup.PUT("/auth/me", h.AuthAvius.UpdateMe)
	}
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
