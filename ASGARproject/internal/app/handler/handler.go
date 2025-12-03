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
	return &Handler{Repository: r}
}

func (h *Handler) RegisterHandler(router *gin.Engine) {
	// ================ ДОМЕН: УСЛУГИ (7 методов) ================
	router.GET("/api/services", h.GetAsgarServices)
	router.GET("/api/services/:id", h.GetAsgarService)
	router.POST("/api/services", h.CreateAsgarService)
	router.PUT("/api/services/:id", h.UpdateAsgarService)
	router.DELETE("/api/services/:id", h.DeleteAsgarService)
	router.POST("/api/services/:id/add-to-miniplane", h.AddServiceToMiniplane)
	router.POST("/api/services/:id/upload", h.UploadServiceImage)

	// ================ ДОМЕН: МИНИПЛЭЙН (2 метода) ================
	router.GET("/api/miniplane/icon", h.GetMiniplaneIcon)
	router.GET("/api/miniplane", h.GetMiniplane)

	// ================ ДОМЕН: ЗАЯВКИ (6 методов) ================
	router.GET("/api/flight-servs", h.GetFlightServs)
	router.GET("/api/flight-servs/:id", h.GetFlightServ)
	router.PUT("/api/flight-servs/:id", h.UpdateFlightServ)
	router.PUT("/api/flight-servs/:id/submit", h.SubmitFlightServ)
	router.PUT("/api/flight-servs/:id/complete", h.CompleteFlightServ)
	router.DELETE("/api/flight-servs/:id", h.DeleteFlightServ)

	// ================ ДОМЕН: M-M СВЯЗИ (2 метода) ================
	router.DELETE("/api/subjservs/:id", h.DeleteSubjserv)
	router.PUT("/api/subjservs/:id", h.UpdateSubjserv)

	// ================ ДОМЕН: АВТОРИЗАЦИЯ (5 методов) ================
	router.POST("/api/auth/register", h.RegisterUser)
	router.GET("/api/users/profile", h.GetUserProfile)
	router.PUT("/api/users/profile", h.UpdateUserProfile)
	router.POST("/api/auth/login", h.LoginUser)
	router.POST("/api/auth/logout", h.LogoutUser)

	// Статические файлы (для загруженных изображений)
	router.Static("/uploads", "./uploads")

	// Для фронтенда (если есть)
	router.Static("/app", "./frontend") // если есть фронтенд
}

func (h *Handler) RegisterStatic(router *gin.Engine) {
	router.LoadHTMLGlob("templates/*")
	router.Static("/static", "./resources")
}

// Общая функция ошибки
func (h *Handler) errorResponse(ctx *gin.Context, statusCode int, err error) {
	logrus.Error(err.Error())
	ctx.JSON(statusCode, gin.H{
		"status": "error",
		"error":  err.Error(),
	})
}
