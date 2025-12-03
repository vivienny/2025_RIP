package handler

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"decode/internal/app/ds"

	"github.com/gin-gonic/gin"
)

// ================ ДОМЕН: ЗАЯВКИ (6 методов) ================

func (h *Handler) GetFlightServs(ctx *gin.Context) {
	statusFilter := ctx.Query("status")
	dateFrom := ctx.Query("date_from")
	dateTo := ctx.Query("date_to")

	requests, err := h.Repository.GetAllRequests()
	if err != nil {
		h.errorResponse(ctx, http.StatusInternalServerError, err)
		return
	}

	var filteredRequests []ds.FlightServ
	for _, req := range requests {
		if req.Status == "удален" || req.Status == "черновик" {
			continue
		}

		if statusFilter != "" && statusFilter != "все" && req.Status != statusFilter {
			continue
		}

		if dateFrom != "" {
			if date, err := time.Parse("2006-01-02", dateFrom); err == nil {
				if req.DateCreate.Before(date) {
					continue
				}
			}
		}
		if dateTo != "" {
			if date, err := time.Parse("2006-01-02", dateTo); err == nil {
				if req.DateCreate.After(date) {
					continue
				}
			}
		}

		filteredRequests = append(filteredRequests, req)
	}

	for i := range filteredRequests {
		requestWithItems, err := h.Repository.GetRequestWithItems(filteredRequests[i].ID)
		if err == nil {
			count := 0
			for _, item := range requestWithItems.Subjservs {
				if strings.TrimSpace(item.Price) != "" && item.Price != "0" {
					count++
				}
			}
			filteredRequests[i].TotalPrice = fmt.Sprintf("%s (items: %d)",
				filteredRequests[i].TotalPrice, count)
		}
	}

	ctx.JSON(http.StatusOK, gin.H{
		"requests": filteredRequests,
		"count":    len(filteredRequests),
		"filters": gin.H{
			"status":    statusFilter,
			"date_from": dateFrom,
			"date_to":   dateTo,
		},
	})
}

func (h *Handler) GetFlightServ(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		h.errorResponse(ctx, http.StatusBadRequest, err)
		return
	}

	request, err := h.Repository.GetRequestWithItems(uint(id))
	if err != nil {
		h.errorResponse(ctx, http.StatusNotFound, fmt.Errorf("request not found"))
		return
	}

	user, err := h.Repository.GetUserByID(request.UserID)
	if err != nil {
		h.errorResponse(ctx, http.StatusInternalServerError, err)
		return
	}

	var itemsWithImages []gin.H
	for _, item := range request.Subjservs {
		itemsWithImages = append(itemsWithImages, gin.H{
			"id":       item.ID,
			"service":  item.Service,
			"quantity": item.Quantity,
			"price":    item.Price,
			"image_url": func() string {
				if item.Service.Img != "" {
					return "/uploads/" + item.Service.Img
				}
				return ""
			}(),
		})
	}

	ctx.JSON(http.StatusOK, gin.H{
		"request": gin.H{
			"id":          request.ID,
			"status":      request.Status,
			"total_price": request.TotalPrice,
			"date_create": request.DateCreate,
			"date_update": request.DateUpdate,
			"date_finish": request.DateFinish,
			"user_login":  user.Login,
			"items":       itemsWithImages,
		},
	})
}

func (h *Handler) UpdateFlightServ(ctx *gin.Context) {
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

	delete(updates, "id")
	delete(updates, "status")
	delete(updates, "user_id")
	delete(updates, "date_create")
	delete(updates, "date_finish")

	if err := h.Repository.UpdateRequest(uint(id), updates); err != nil {
		h.errorResponse(ctx, http.StatusInternalServerError, err)
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "Request updated",
	})
}

func (h *Handler) SubmitFlightServ(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		h.errorResponse(ctx, http.StatusBadRequest, err)
		return
	}

	request, err := h.Repository.GetRequestWithItems(uint(id))
	if err != nil {
		h.errorResponse(ctx, http.StatusNotFound, fmt.Errorf("request not found"))
		return
	}

	if request.Status != "черновик" {
		h.errorResponse(ctx, http.StatusBadRequest,
			fmt.Errorf("only draft requests can be submitted"))
		return
	}

	if len(request.Subjservs) == 0 {
		h.errorResponse(ctx, http.StatusBadRequest,
			fmt.Errorf("request must have at least one service"))
		return
	}

	if err := h.Repository.UpdateRequestStatus(uint(id), "формирование"); err != nil {
		h.errorResponse(ctx, http.StatusInternalServerError, err)
		return
	}

	if err := h.Repository.UpdateRequest(uint(id), map[string]interface{}{
		"date_update": time.Now(),
	}); err != nil {
		h.errorResponse(ctx, http.StatusInternalServerError, err)
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message":    "Request submitted",
		"new_status": "формирование",
	})
}

func (h *Handler) CompleteFlightServ(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		h.errorResponse(ctx, http.StatusBadRequest, err)
		return
	}

	action := ctx.Query("action")
	if action == "" {
		action = "complete"
	}

	request, err := h.Repository.GetRequestWithItems(uint(id))
	if err != nil {
		h.errorResponse(ctx, http.StatusNotFound, fmt.Errorf("request not found"))
		return
	}

	if request.Status != "формирование" {
		h.errorResponse(ctx, http.StatusBadRequest,
			fmt.Errorf("only requests in 'формирование' can be completed"))
		return
	}

	var totalCost float64
	var deliveryCost float64 = 500

	for _, item := range request.Subjservs {
		cleanPrice := strings.ReplaceAll(item.Price, " ", "")
		price, _ := strconv.ParseFloat(cleanPrice, 64)
		totalCost += price * float64(item.Quantity)
	}

	finalCost := totalCost + deliveryCost
	newStatus := "завершено"
	if action == "reject" {
		newStatus = "отклонено"
	}

	updates := map[string]interface{}{
		"status":      newStatus,
		"date_finish": time.Now(),
		"date_update": time.Now(),
		"total_price": fmt.Sprintf("%.2f", finalCost),
	}

	if err := h.Repository.UpdateRequest(uint(id), updates); err != nil {
		h.errorResponse(ctx, http.StatusInternalServerError, err)
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": fmt.Sprintf("Request %s", newStatus),
		"calculation": gin.H{
			"services_cost":  totalCost,
			"delivery_cost":  deliveryCost,
			"final_cost":     finalCost,
			"services_count": len(request.Subjservs),
			"action":         action,
		},
	})
}

func (h *Handler) DeleteFlightServ(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		h.errorResponse(ctx, http.StatusBadRequest, err)
		return
	}

	if err := h.Repository.UpdateRequestStatus(uint(id), "удален"); err != nil {
		h.errorResponse(ctx, http.StatusInternalServerError, err)
		return
	}

	if err := h.Repository.UpdateRequest(uint(id), map[string]interface{}{
		"date_update": time.Now(),
	}); err != nil {
		h.errorResponse(ctx, http.StatusInternalServerError, err)
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message":    "Request deleted",
		"new_status": "удален",
	})
}
