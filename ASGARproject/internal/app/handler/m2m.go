package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// ================ ДОМЕН: M-M СВЯЗИ (2 метода) ================

func (h *Handler) DeleteSubjserv(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		h.errorResponse(ctx, http.StatusBadRequest, err)
		return
	}

	if err := h.Repository.DeleteRequestItem(uint(id)); err != nil {
		h.errorResponse(ctx, http.StatusInternalServerError, err)
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "Position removed",
	})
}

func (h *Handler) UpdateSubjserv(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		h.errorResponse(ctx, http.StatusBadRequest, err)
		return
	}

	var update struct {
		Quantity int `json:"quantity" binding:"min=1"`
	}

	if err := ctx.BindJSON(&update); err != nil {
		h.errorResponse(ctx, http.StatusBadRequest, err)
		return
	}

	if err := h.Repository.UpdateRequestItem(uint(id), update.Quantity); err != nil {
		h.errorResponse(ctx, http.StatusInternalServerError, err)
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message":      "Position updated",
		"new_quantity": update.Quantity,
	})
}
