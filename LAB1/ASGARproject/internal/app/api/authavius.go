package api

import (
	"net/http"

	"decode/internal/app/ds"
	"decode/internal/app/repository"

	"github.com/gin-gonic/gin"
)

type AuthAviusHandler struct {
	Repository *repository.Repository
}

func NewAuthAviusHandler(r *repository.Repository) *AuthAviusHandler {
	return &AuthAviusHandler{
		Repository: r,
	}
}

// Register регистрация пользователя
// POST /api/auth/register
func (h *AuthAviusHandler) Register(ctx *gin.Context) {
	var user ds.Avius

	if err := ctx.BindJSON(&user); err != nil {
		ctx.JSON(http.StatusBadRequest, nil)
		return
	}

	// Проверяем нет ли пользователя с таким логином
	_, err := h.Repository.GetUserByLogin(user.Login)
	if err == nil {
		ctx.JSON(http.StatusConflict, nil)
		return
	}

	// Создаем пользователя
	err = h.Repository.CreateUser(&user)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, nil)
		return
	}

	ctx.JSON(http.StatusCreated, user)
}

// Login вход пользователя
// POST /api/auth/login
func (h *AuthAviusHandler) Login(ctx *gin.Context) {
	var request struct {
		Login    string `json:"login"`
		Password string `json:"password"`
	}

	if err := ctx.BindJSON(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, nil)
		return
	}

	// Находим пользователя
	user, err := h.Repository.GetUserByLogin(request.Login)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, nil)
		return
	}

	// Проверяем пароль (пока без хеширования)
	if user.Password != request.Password {
		ctx.JSON(http.StatusUnauthorized, nil)
		return
	}

	ctx.JSON(http.StatusOK, user)
}

// Logout выход пользователя
// POST /api/auth/logout
func (h *AuthAviusHandler) Logout(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, nil)
}

// GetMe получение данных пользователя
// GET /api/auth/me
func (h *AuthAviusHandler) GetMe(ctx *gin.Context) {
	// TODO: Получать ID из сессии/JWT
	userID := uint(1) // временно

	user, err := h.Repository.GetUserByID(userID)
	if err != nil {
		ctx.JSON(http.StatusNotFound, nil)
		return
	}

	ctx.JSON(http.StatusOK, user)
}

// UpdateMe обновление данных пользователя
// PUT /api/auth/me
func (h *AuthAviusHandler) UpdateMe(ctx *gin.Context) {
	// TODO: Получать ID из сессии/JWT
	userID := uint(1) // временно

	var updates map[string]interface{}
	if err := ctx.BindJSON(&updates); err != nil {
		ctx.JSON(http.StatusBadRequest, nil)
		return
	}

	user, err := h.Repository.GetUserByID(userID)
	if err != nil {
		ctx.JSON(http.StatusNotFound, nil)
		return
	}

	// Обновляем поля
	if login, ok := updates["Login"].(string); ok {
		user.Login = login
	}
	if password, ok := updates["Password"].(string); ok {
		user.Password = password
	}

	err = h.Repository.UpdateUser(&user)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, nil)
		return
	}

	ctx.JSON(http.StatusOK, user)
}
