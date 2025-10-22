package api

import (
	"net/http"
	"strconv"

	"decode/internal/app/repository"

	"github.com/gin-gonic/gin"
)

type MiniplaneHandler struct {
	Repository *repository.Repository
}

func NewMiniplaneHandler(r *repository.Repository) *MiniplaneHandler {
	return &MiniplaneHandler{
		Repository: r,
	}
}

// MiniplaneItemResponse ответ для позиции корзины
type MiniplaneItemResponse struct {
	ID        uint   `json:"id"`
	ServiceID uint   `json:"service_id"`
	Name      string `json:"name"`
	Price     string `json:"price"`
	Quantity  int    `json:"quantity"`
}

// MiniplaneResponse ответ для корзины
type MiniplaneResponse struct {
	ID    uint                    `json:"id"`
	Items []MiniplaneItemResponse `json:"items"`
	Total int                     `json:"total"`
	Count int                     `json:"count"`
}

// GetMiniplaneIconAPI возвращает иконку корзины (id + количество)
// GET /api/miniplane/icon
func (h *MiniplaneHandler) GetMiniplaneIconAPI(ctx *gin.Context) {
	userID := uint(1)

	count := h.Repository.GetMiniplaneItemsCount()

	// Получаем ID активной корзины
	miniplaneID, err := h.Repository.GetActiveMiniplaneID(userID)
	if err != nil {
		ctx.JSON(http.StatusOK, gin.H{
			"miniplane_id": 0,
			"count":        0,
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"miniplane_id": miniplaneID,
		"count":        count,
	})
}

// GetMiniplaneAPI возвращает корзину в JSON формате
// GET /api/miniplane
func (h *MiniplaneHandler) GetMiniplaneAPI(ctx *gin.Context) {
	userID := uint(1)

	// Получаем активную корзину
	miniplaneID, err := h.Repository.GetActiveMiniplaneID(userID)
	if err != nil {
		ctx.JSON(http.StatusOK, MiniplaneResponse{
			ID:    0,
			Items: []MiniplaneItemResponse{},
			Total: 0,
			Count: 0,
		})
		return
	}

	// Получаем позиции корзины
	items, err := h.Repository.GetMiniplaneItems(miniplaneID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, nil)
		return
	}

	// Получаем общую сумму
	total, err := h.Repository.GetMiniplaneTotal(miniplaneID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, nil)
		return
	}

	// Формируем ответ
	responseItems := make([]MiniplaneItemResponse, len(items))
	for i, item := range items {
		responseItems[i] = MiniplaneItemResponse{
			ID:        item.ID,
			ServiceID: item.ServiceID,
			Name:      item.Service.Name,
			Price:     item.Price,
			Quantity:  item.Quantity,
		}
	}

	response := MiniplaneResponse{
		ID:    miniplaneID,
		Items: responseItems,
		Total: total,
		Count: len(items),
	}

	ctx.JSON(http.StatusOK, response)
}

// RemoveFromMiniplaneAPI удаляет позицию из корзины (API)
// DELETE /api/miniplane/items/:id
func (h *MiniplaneHandler) RemoveFromMiniplaneAPI(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil || id <= 0 {
		ctx.JSON(http.StatusBadRequest, nil)
		return
	}

	err = h.Repository.RemoveFromMiniplane(uint(id))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, nil)
		return
	}

	ctx.JSON(http.StatusOK, nil)
}

// UpdateMiniplaneItemAPI изменяет количество в позиции корзины
// PUT /api/miniplane/items/:id
func (h *MiniplaneHandler) UpdateMiniplaneItemAPI(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil || id <= 0 {
		ctx.JSON(http.StatusBadRequest, nil)
		return
	}

	var request struct {
		Quantity int `json:"quantity"`
	}

	if err := ctx.BindJSON(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, nil)
		return
	}

	if request.Quantity <= 0 {
		ctx.JSON(http.StatusBadRequest, nil)
		return
	}

	ctx.JSON(http.StatusOK, nil)
}

// ClearMiniplaneAPI очищает всю корзину (API)
// DELETE /api/miniplane/clear
func (h *MiniplaneHandler) ClearMiniplaneAPI(ctx *gin.Context) {
	userID := uint(1)

	// Получаем активную корзину
	miniplaneID, err := h.Repository.GetActiveMiniplaneID(userID)
	if err != nil {
		ctx.JSON(http.StatusOK, nil)
		return
	}

	err = h.Repository.ClearMiniplane(miniplaneID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, nil)
		return
	}

	ctx.JSON(http.StatusOK, nil)
}

func (h *MiniplaneHandler) GetMiniplaneCountAPI(ctx *gin.Context) {
	count := h.Repository.GetMiniplaneItemsCount()

	ctx.JSON(http.StatusOK, gin.H{
		"count": count,
	})
}
