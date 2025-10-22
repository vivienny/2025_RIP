package api

import (
	"net/http"
	"strconv"

	"decode/internal/app/ds"
	"decode/internal/app/repository"

	"github.com/gin-gonic/gin"
)

type ServicesHandler struct {
	Repository *repository.Repository
}

func NewServicesHandler(r *repository.Repository) *ServicesHandler {
	return &ServicesHandler{
		Repository: r,
	}
}

// GetServices возвращает список услуг с фильтрацией
// GET /api/services
func (h *ServicesHandler) GetServices(ctx *gin.Context) {
	var services []ds.ASGARService
	var err error

	// Фильтрация по имени (если передана)
	search := ctx.Query("search")
	if search == "" {
		services, err = h.Repository.GetAllServices()
	} else {
		services, err = h.Repository.SearchServicesByName(search)
	}

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, nil)
		return
	}

	ctx.JSON(http.StatusOK, services)
}

// GetService возвращает одну услугу по ID
// GET /api/services/:id
func (h *ServicesHandler) GetService(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil || id <= 0 {
		ctx.JSON(http.StatusBadRequest, nil)
		return
	}

	service, err := h.Repository.GetAsgarService(id)
	if err != nil {
		ctx.JSON(http.StatusNotFound, nil)
		return
	}

	ctx.JSON(http.StatusOK, service)
}

// CreateService создает новую услугу
// POST /api/services
func (h *ServicesHandler) CreateService(ctx *gin.Context) {
	var service ds.ASGARService

	if err := ctx.BindJSON(&service); err != nil {
		ctx.JSON(http.StatusBadRequest, nil)
		return
	}

	// Создаем услугу в базе
	err := h.Repository.CreateService(&service)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, nil)
		return
	}

	ctx.JSON(http.StatusCreated, service)
}

// UpdateService обновляет услугу
// PUT /api/services/:id
func (h *ServicesHandler) UpdateService(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil || id <= 0 {
		ctx.JSON(http.StatusBadRequest, nil)
		return
	}

	var updates map[string]interface{}
	if err := ctx.BindJSON(&updates); err != nil {
		ctx.JSON(http.StatusBadRequest, nil)
		return
	}

	// Обновляем услугу
	err = h.Repository.UpdateService(uint(id), updates)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, nil)
		return
	}

	ctx.JSON(http.StatusOK, nil)
}

// DeleteService удаляет услугу (мягкое удаление)
// DELETE /api/services/:id
func (h *ServicesHandler) DeleteService(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil || id <= 0 {
		ctx.JSON(http.StatusBadRequest, nil)
		return
	}

	// Удаляем услугу
	err = h.Repository.DeleteService(uint(id))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, nil)
		return
	}

	ctx.JSON(http.StatusOK, nil)
}

// UploadServiceImage загружает изображение для услуги
// POST /api/services/:id/image
func (h *ServicesHandler) UploadServiceImage(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil || id <= 0 {
		ctx.JSON(http.StatusBadRequest, nil)
		return
	}

	ctx.JSON(http.StatusOK, nil)
}
