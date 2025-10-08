package handler

import (
	"decode/internal/app/repository"
	"net/http"
	"strconv"

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

func (h *Handler) GetAsgarList(ctx *gin.Context) {
	var aviationServices []repository.ASGARserv
	var err error

	findavia := ctx.Query("query")
	if findavia == "" {
		aviationServices, err = h.Repository.GetAsgarList()
		if err != nil {
			logrus.Error(err)
		}
	} else {
		aviationServices, err = h.Repository.GetFindAviaSub(findavia)
		if err != nil {
			logrus.Error(err)
		}
	}

	ctx.HTML(http.StatusOK, "asgar_catalog.html", gin.H{
		"aviationServices": aviationServices,
		"findavia":         findavia,
	})
}

func (h *Handler) GetSelectAviaSub(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		logrus.Error(err)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Неверный ID"})
		return
	}

	selsubavia, err := h.Repository.GetSelectAviaSub(id)
	if err != nil {
		logrus.Error(err)
		ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
	}

	ctx.HTML(http.StatusOK, "selected_aviaproduct.html", gin.H{
		"selsubavia": selsubavia,
	})
}

func (h *Handler) GetMiniPlane(ctx *gin.Context) {
	ctx.HTML(http.StatusOK, "miniplane.html", gin.H{})
}
