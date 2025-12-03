package handler

import (
	"fmt"
	"net/http"

	"decode/internal/app/ds"

	"github.com/gin-gonic/gin"
)

// ================ ДОМЕН: АВТОРИЗАЦИЯ (5 методов) ================

func (h *Handler) RegisterUser(ctx *gin.Context) {
	var request struct {
		Login    string `json:"login" binding:"required,min=3"`
		Password string `json:"password" binding:"required,min=6"`
	}

	if err := ctx.BindJSON(&request); err != nil {
		h.errorResponse(ctx, http.StatusBadRequest, err)
		return
	}

	var existingUser ds.Avius
	if err := h.Repository.GetUserByLogin(request.Login, &existingUser); err == nil {
		h.errorResponse(ctx, http.StatusBadRequest, fmt.Errorf("user already exists"))
		return
	}

	user := ds.Avius{
		Login:       request.Login,
		Password:    request.Password,
		IsModerator: false,
	}

	if err := h.Repository.CreateUser(&user); err != nil {
		h.errorResponse(ctx, http.StatusInternalServerError, err)
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{
		"message": "User registered",
		"user_id": user.ID,
	})
}

func (h *Handler) GetUserProfile(ctx *gin.Context) {
	userID := GetCurrentUserID()
	user, err := h.Repository.GetUserByID(userID)
	if err != nil {
		h.errorResponse(ctx, http.StatusNotFound, fmt.Errorf("user not found"))
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"user": gin.H{
			"id":           user.ID,
			"login":        user.Login,
			"is_moderator": user.IsModerator,
		},
	})
}

func (h *Handler) UpdateUserProfile(ctx *gin.Context) {
	userID := GetCurrentUserID()
	var updates map[string]interface{}
	if err := ctx.BindJSON(&updates); err != nil {
		h.errorResponse(ctx, http.StatusBadRequest, err)
		return
	}

	allowedFields := map[string]bool{"login": true, "password": true}
	filteredUpdates := make(map[string]interface{})
	for key, value := range updates {
		if allowedFields[key] {
			filteredUpdates[key] = value
		}
	}

	if len(filteredUpdates) == 0 {
		h.errorResponse(ctx, http.StatusBadRequest, fmt.Errorf("no valid fields"))
		return
	}

	if err := h.Repository.UpdateUser(userID, filteredUpdates); err != nil {
		h.errorResponse(ctx, http.StatusInternalServerError, err)
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "Profile updated",
		"user_id": userID,
	})
}

func (h *Handler) LoginUser(ctx *gin.Context) {
	var request struct {
		Login    string `json:"login" binding:"required"`
		Password string `json:"password" binding:"required"`
	}

	if err := ctx.BindJSON(&request); err != nil {
		h.errorResponse(ctx, http.StatusBadRequest, err)
		return
	}

	var user ds.Avius
	if err := h.Repository.GetUserByLogin(request.Login, &user); err != nil {
		h.errorResponse(ctx, http.StatusUnauthorized, fmt.Errorf("invalid credentials"))
		return
	}

	if user.Password != request.Password {
		h.errorResponse(ctx, http.StatusUnauthorized, fmt.Errorf("invalid credentials"))
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message":      "Login successful",
		"user_id":      user.ID,
		"is_moderator": user.IsModerator,
	})
}

func (h *Handler) LogoutUser(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, gin.H{
		"message": "Logout successful",
	})
}
