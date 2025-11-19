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

// RegisterHandler - ТЕ ЖЕ ПУТИ ЧТО В ЛР2, НО ВОЗВРАЩАЮТ JSON
func (h *Handler) RegisterHandler(router *gin.Engine) {
	// ⭐ ТОЧНО ТЕ ЖЕ ПУТИ КАК В ЛР2
	router.GET("/", h.GetAsgarServices)
	router.GET("/service/:id", h.GetAsgarService)
	router.GET("/miniplane", h.GetMiniplane)
	router.POST("/remove-from-miniplane", h.RemoveFromMiniplane)
	router.POST("/clear-miniplane", h.ClearMiniplane)
	router.POST("/add-to-cart", h.AddToCart)

	router.POST("/create-aviasub", h.CreateService)
	router.PUT("/service/:id", h.UpdateService)

	router.DELETE("/service/:id", h.DeleteMiniplane)
	router.GET("/requests/miniplane-icon", h.GetMiniplaneIcon)

	router.GET("/requests", h.GetRequests)
	router.GET("/requests/:id", h.GetRequest)

	router.PUT("/requests/:id", h.UpdateRequest)
	router.PUT("/requests/:id/form", h.FormRequest)
	router.PUT("/requests/:id/complete", h.CompleteRequest)

	router.DELETE("/request-items/:id", h.DeleteRequestItem)
	router.PUT("/request-items/:id", h.UpdateRequestItem)

	router.POST("/auth/register", h.RegisterUser)
	router.GET("/users/profile", h.GetUserProfile)
	router.PUT("/users/profile", h.UpdateUserProfile)

	router.POST("/auth/login", h.LoginUser)
	router.POST("/auth/logout", h.LogoutUser)
}

// RegisterStatic оставляем как было
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
