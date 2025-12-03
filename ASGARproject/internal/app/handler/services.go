package handler

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"decode/internal/app/ds"

	"github.com/gin-gonic/gin"
)

// ================ ДОМЕН: УСЛУГИ (7 методов) ================

func (h *Handler) GetAsgarServices(ctx *gin.Context) {
	search := ctx.Query("search")

	// ВСЕГДА используем SearchServicesByName, даже если search пустой
	services, err := h.Repository.SearchServicesByName(search)

	if err != nil {
		h.errorResponse(ctx, http.StatusInternalServerError, err)
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"services":   services,
		"count":      len(services),
		"search":     search,
		"cart_count": h.Repository.GetCartItemsCount(),
	})
}

func (h *Handler) GetAsgarService(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		h.errorResponse(ctx, http.StatusBadRequest, err)
		return
	}

	service, err := h.Repository.GetAsgarService(id)
	if err != nil {
		h.errorResponse(ctx, http.StatusNotFound, err)
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"service": service,
	})
}

func (h *Handler) CreateAsgarService(ctx *gin.Context) {
	var service ds.ASGARService
	if err := ctx.BindJSON(&service); err != nil {
		h.errorResponse(ctx, http.StatusBadRequest, err)
		return
	}

	service.IsDelete = false
	service.Img = ""

	if err := h.Repository.CreateService(&service); err != nil {
		h.errorResponse(ctx, http.StatusInternalServerError, err)
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{
		"service": service,
	})
}

func (h *Handler) UpdateAsgarService(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		h.errorResponse(ctx, http.StatusBadRequest, err)
		return
	}

	var updates map[string]interface{}
	if err := ctx.BindJSON(&updates); err != nil {
		h.errorResponse(ctx, http.StatusBadRequest, err)
		return
	}

	allowedFields := map[string]bool{
		"name": true, "info": true, "price": true,
		"full_description": true, "unit": true,
	}

	filteredUpdates := make(map[string]interface{})
	for key, value := range updates {
		if allowedFields[key] {
			filteredUpdates[key] = value
		}
	}

	delete(filteredUpdates, "id")
	delete(filteredUpdates, "is_delete")
	delete(filteredUpdates, "img")

	if len(filteredUpdates) == 0 {
		h.errorResponse(ctx, http.StatusBadRequest, fmt.Errorf("no valid fields"))
		return
	}

	if err := h.Repository.UpdateService(uint(id), filteredUpdates); err != nil {
		h.errorResponse(ctx, http.StatusInternalServerError, err)
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "Service updated",
	})
}

func (h *Handler) DeleteAsgarService(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		h.errorResponse(ctx, http.StatusBadRequest, err)
		return
	}

	if err := h.Repository.UpdateService(uint(id), map[string]interface{}{
		"is_delete": true,
	}); err != nil {
		h.errorResponse(ctx, http.StatusInternalServerError, err)
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "Service deleted",
	})
}

func (h *Handler) AddServiceToMiniplane(ctx *gin.Context) {
	userID := GetCurrentUserID()
	idStr := ctx.Param("id")
	serviceID, err := strconv.Atoi(idStr)
	if err != nil {
		h.errorResponse(ctx, http.StatusBadRequest, err)
		return
	}

	flightServID, err := h.Repository.GetOrCreateCartFlightServ(userID)
	if err != nil {
		h.errorResponse(ctx, http.StatusInternalServerError, err)
		return
	}

	err = h.Repository.AddToCart(flightServID, uint(serviceID), 1)
	if err != nil {
		h.errorResponse(ctx, http.StatusInternalServerError, err)
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message":      "Service added to miniplane",
		"miniplane_id": flightServID,
		"service_id":   serviceID,
	})
}

func (h *Handler) UploadServiceImage(ctx *gin.Context) {
	idStr := ctx.Param("id")
	serviceID, err := strconv.Atoi(idStr)
	if err != nil {
		h.errorResponse(ctx, http.StatusBadRequest, err)
		return
	}

	file, err := ctx.FormFile("image")
	if err != nil {
		h.errorResponse(ctx, http.StatusBadRequest, fmt.Errorf("no image file"))
		return
	}

	filename := fmt.Sprintf("service_%d_%s", serviceID,
		strings.ToLower(strings.ReplaceAll(file.Filename, " ", "_")))

	uploadPath := "./uploads/" + filename
	if err := ctx.SaveUploadedFile(file, uploadPath); err != nil {
		h.errorResponse(ctx, http.StatusInternalServerError, err)
		return
	}

	if err := h.Repository.UpdateService(uint(serviceID), map[string]interface{}{
		"img": filename,
	}); err != nil {
		h.errorResponse(ctx, http.StatusInternalServerError, err)
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message":  "Image uploaded",
		"filename": filename,
		"url":      "/uploads/" + filename,
	})
}
