package handler

import (
	"net/http"
	"strconv"
	"strings"

	"decode/internal/app/ds"

	"github.com/gin-gonic/gin"
)

// ================ ДОМЕН: МИНИПЛЭЙН (2 метода) ================

func (h *Handler) GetMiniplaneIcon(ctx *gin.Context) {
	userID := GetCurrentUserID()
	flightServID, err := h.Repository.GetOrCreateCartFlightServ(userID)
	if err != nil {
		ctx.JSON(http.StatusOK, gin.H{
			"request_id": 0,
			"item_count": 0,
			"total_sum":  0,
		})
		return
	}

	items, err := h.Repository.GetCartItems(flightServID)
	if err != nil {
		h.errorResponse(ctx, http.StatusInternalServerError, err)
		return
	}

	total := 0
	for _, item := range items {
		cleanPrice := strings.ReplaceAll(item.Price, " ", "")
		price, _ := strconv.Atoi(cleanPrice)
		total += price * item.Quantity
	}

	ctx.JSON(http.StatusOK, gin.H{
		"request_id": flightServID,
		"item_count": len(items),
		"total_sum":  total,
	})
}

func (h *Handler) GetMiniplane(ctx *gin.Context) {
	userID := GetCurrentUserID()
	flightServID, err := h.Repository.GetOrCreateCartFlightServ(userID)
	if err != nil {
		ctx.JSON(http.StatusOK, gin.H{
			"items": []ds.Subjserv{},
			"total": 0,
		})
		return
	}

	items, err := h.Repository.GetCartItems(flightServID)
	if err != nil {
		h.errorResponse(ctx, http.StatusInternalServerError, err)
		return
	}

	total := 0
	for _, item := range items {
		cleanPrice := strings.ReplaceAll(item.Price, " ", "")
		price, err := strconv.Atoi(cleanPrice)
		if err != nil {
			continue
		}
		total += price * item.Quantity
	}

	ctx.JSON(http.StatusOK, gin.H{
		"miniplane_id": flightServID,
		"items":        items,
		"total":        total,
		"item_count":   len(items),
	})
}
