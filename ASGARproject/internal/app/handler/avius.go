package handler

import (
	"fmt"
	"net/http"
	"time"

	"decode/internal/app/ds"

	"crypto/sha1"
	"encoding/hex"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

// Структуры для Swagger документации
type RegisterRequest struct {
	Login    string `json:"login" binding:"required,min=3" example:"user123"`
	Password string `json:"password" binding:"required,min=6" example:"password123"`
}

type LoginRequest struct {
	Login    string `json:"login" binding:"required" example:"user123"`
	Password string `json:"password" binding:"required" example:"password123"`
}

// JWT секретный ключ (для демонстрации, в продакшене брать из конфига)
const jwtSecret = "asgar-avia-secret-key-2025"

// ================ ДОМЕН: АВТОРИЗАЦИЯ (5 методов) ================

// RegisterUser godoc
// @Summary Регистрация пользователя
// @Description Регистрация нового пользователя в системе с SHA-1 хешированием пароля
// @Tags auth
// @Accept json
// @Produce json
// @Param request body RegisterRequest true "Данные для регистрации"
// @Success 201 {object} map[string]interface{} "Успешная регистрация"
// @Failure 400 {object} map[string]interface{} "Неверные данные"
// @Failure 500 {object} map[string]interface{} "Внутренняя ошибка сервера"
// @Router /auth/register [post]
func (h *Handler) RegisterUser(ctx *gin.Context) {
	var request RegisterRequest

	if err := ctx.BindJSON(&request); err != nil {
		h.errorResponse(ctx, http.StatusBadRequest, err)
		return
	}

	existingUser, err := h.Repository.GetUserByLogin(request.Login)
	if err == nil && existingUser != nil {
		h.errorResponse(ctx, http.StatusBadRequest, fmt.Errorf("user already exists"))
		return
	}

	// Хешируем пароль как в методичке (SHA1)
	hashedPassword := generateHashString(request.Password)

	user := ds.Avius{
		Login:       request.Login,
		Password:    hashedPassword, // ← ХЕШИРОВАННЫЙ пароль
		IsModerator: false,
		Role:        "user",
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

// Функция из методички
func generateHashString(s string) string {
	h := sha1.New()
	h.Write([]byte(s))
	return hex.EncodeToString(h.Sum(nil))
}

// GetUserProfile godoc
// @Summary Получение профиля пользователя
// @Description Получение информации о текущем пользователе
// @Tags users
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} map[string]interface{} "Успешное получение профиля"
// @Failure 401 {object} map[string]interface{} "Не авторизован"
// @Failure 404 {object} map[string]interface{} "Пользователь не найден"
// @Router /users/profile [get]
func (h *Handler) GetUserProfile(ctx *gin.Context) {
	// Временное решение для компиляции
	userID := uint(1) // ← Вместо GetCurrentUserID()

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

// UpdateUserProfile godoc
// @Summary Обновление профиля пользователя
// @Description Обновление данных пользователя (логин, пароль)
// @Tags users
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param updates body map[string]interface{} true "Обновляемые поля"
// @Success 200 {object} map[string]interface{} "Профиль обновлен"
// @Failure 400 {object} map[string]interface{} "Неверные данные"
// @Failure 401 {object} map[string]interface{} "Не авторизован"
// @Router /users/profile [put]
func (h *Handler) UpdateUserProfile(ctx *gin.Context) {
	// Временное решение для компиляции
	userID := uint(1)
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

// LoginUser godoc
// @Summary Аутентификация пользователя
// @Description Вход пользователя с проверкой SHA-1 хеша пароля и выдачей JWT токена
// @Tags auth
// @Accept json
// @Produce json
// @Param request body LoginRequest true "Учетные данные"
// @Success 200 {object} map[string]interface{} "Успешный вход"
// @Failure 400 {object} map[string]interface{} "Неверный запрос"
// @Failure 401 {object} map[string]interface{} "Неверные учетные данные"
// @Router /auth/login [post]
func (h *Handler) LoginUser(ctx *gin.Context) {
	var request LoginRequest

	if err := ctx.BindJSON(&request); err != nil {
		h.errorResponse(ctx, http.StatusBadRequest, err)
		return
	}

	user, err := h.Repository.GetUserByLogin(request.Login)
	if err != nil {
		h.errorResponse(ctx, http.StatusUnauthorized, fmt.Errorf("invalid credentials"))
		return
	}

	// Проверяем хеш SHA1 (как в методичке)
	hashedInput := generateHashString(request.Password)

	if user.Password != hashedInput {
		h.errorResponse(ctx, http.StatusUnauthorized, fmt.Errorf("invalid credentials"))
		return
	}

	// ==================== ГЕНЕРАЦИЯ JWT ТОКЕНА ====================
	// Используем ds.JWTClaims как в middleware!
	claims := &ds.JWTClaims{
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Issuer:    "asgar-avia",
		},
		UserID: user.ID,
		Login:  user.Login,
		Role:   user.Role,
	}

	// Создаем токен с алгоритмом HS256
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Подписываем токен секретным ключом
	// ВАЖНО: используйте тот же секрет что в middleware!
	signedToken, err := token.SignedString([]byte(h.JWTConfig.Secret))
	if err != nil {
		h.errorResponse(ctx, http.StatusInternalServerError,
			fmt.Errorf("could not generate token: %v", err))
		return
	}
	// ==================== КОНЕЦ ГЕНЕРАЦИИ JWT ====================

	ctx.JSON(http.StatusOK, gin.H{
		"message":      "Login successful",
		"user_id":      user.ID,
		"login":        user.Login,
		"is_moderator": user.IsModerator,
		"role":         user.Role,
		"token":        signedToken, // ← ДОБАВЛЯЕМ JWT ТОКЕН!
		"token_type":   "Bearer",
		"expires_in":   86400, // 24 часа в секундах
	})
}

// LogoutUser godoc
// @Summary Выход пользователя
// @Description Завершение сессии пользователя (токен добавляется в Redis blacklist)
// @Tags auth
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} map[string]interface{} "Успешный выход"
// @Failure 401 {object} map[string]interface{} "Не авторизован"
// @Router /auth/logout [post]
func (h *Handler) LogoutUser(ctx *gin.Context) {
	// Проверяем, есть ли Redis клиент
	if h.Redis == nil {
		ctx.JSON(http.StatusOK, gin.H{
			"message": "Logout successful (Redis not configured)",
		})
		return
	}

	// Получаем токен из заголовка
	authHeader := ctx.GetHeader("Authorization")
	if authHeader == "" || len(authHeader) < 8 {
		ctx.JSON(http.StatusOK, gin.H{
			"message": "Logout successful (no token provided)",
		})
		return
	}

	// Проверяем префикс Bearer
	if len(authHeader) > 7 && authHeader[:7] == "Bearer " {
		tokenString := authHeader[7:]

		// Добавляем токен в Redis blacklist на 24 часа
		err := h.Redis.AddToBlacklist(ctx, tokenString, 24*time.Hour)
		if err != nil {
			// Логируем ошибку, но все равно считаем logout успешным
			fmt.Printf("Failed to add token to Redis blacklist: %v\n", err)
			ctx.JSON(http.StatusOK, gin.H{
				"message": "Logout successful (but failed to add to blacklist)",
				"error":   err.Error(),
			})
			return
		}

		// Для демонстрации: показываем короткую версию токена
		tokenPreview := tokenString
		if len(tokenString) > 20 {
			tokenPreview = tokenString[:20] + "..."
		}

		ctx.JSON(http.StatusOK, gin.H{
			"message":       "Logout successful",
			"action":        "JWT token added to Redis blacklist",
			"token_preview": tokenPreview,
			"ttl":           "24 hours",
			"redis_key":     "asgar:blacklist:" + tokenString,
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "Logout successful (invalid authorization header)",
	})
}
