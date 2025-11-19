package handler

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"decode/internal/app/ds"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

func (h *Handler) GetAsgarServices(ctx *gin.Context) {
	var services []ds.ASGARService
	var err error

	search := ctx.Query("search")
	if search == "" {
		services, err = h.Repository.GetAllServices()
	} else {
		services, err = h.Repository.SearchServicesByName(search)
	}

	if err != nil {
		h.errorResponse(ctx, http.StatusInternalServerError, err)
		return
	}

	// ✅ JSON вместо HTML
	ctx.JSON(http.StatusOK, gin.H{

		"services":   services,
		"count":      len(services),
		"search":     search,
		"cart_count": h.Repository.GetCartItemsCount(),
	})
}

func (h *Handler) GetMiniplane(ctx *gin.Context) {
	userID := uint(1)

	logrus.Printf("🛒 Получение корзины для UserID=%d", userID)

	flightServID, err := h.Repository.GetOrCreateCartFlightServ(userID)
	if err != nil {
		logrus.Printf("❌ Корзина не найдена для UserID=%d", userID)
		// ✅ JSON вместо HTML
		ctx.JSON(http.StatusOK, gin.H{

			"items": []ds.Subjserv{},
			"total": 0,
		})
		return
	}

	items, err := h.Repository.GetCartItems(flightServID)
	if err != nil {
		logrus.Printf("❌ Ошибка получения позиций корзины: %v", err)
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

	// ✅ JSON вместо HTML
	ctx.JSON(http.StatusOK, gin.H{

		"items": items,
		"total": total,
	})
}

func (h *Handler) RemoveFromMiniplane(ctx *gin.Context) {
	// ✅ JSON вместо Form data
	var request struct {
		SubjservID uint `json:"subjserv_id" binding:"required"`
	}

	if err := ctx.BindJSON(&request); err != nil {
		h.errorResponse(ctx, http.StatusBadRequest, err)
		return
	}

	err := h.Repository.RemoveFromCart(request.SubjservID)
	if err != nil {
		h.errorResponse(ctx, http.StatusInternalServerError, err)
		return
	}

	// ✅ JSON вместо Redirect
	ctx.JSON(http.StatusOK, gin.H{

		"message": "Item removed from cart",
	})
}

func (h *Handler) ClearMiniplane(ctx *gin.Context) {
	userID := uint(1)

	flightServID, err := h.Repository.GetOrCreateCartFlightServ(userID)
	if err != nil {
		// ✅ JSON вместо Redirect
		ctx.JSON(http.StatusOK, gin.H{

			"message": "Cart is already empty",
		})
		return
	}

	err = h.Repository.ClearCart(flightServID)
	if err != nil {
		h.errorResponse(ctx, http.StatusInternalServerError, err)
		return
	}

	// ✅ JSON вместо Redirect
	ctx.JSON(http.StatusOK, gin.H{

		"message": "Cart cleared",
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

	// ✅ JSON вместо HTML
	ctx.JSON(http.StatusOK, gin.H{

		"service": service,
	})
}

func (h *Handler) AddToCart(ctx *gin.Context) {
	userID := uint(1)

	// ✅ JSON вместо Form data
	var request struct {
		ServiceID uint `json:"service_id" binding:"required"`
		Quantity  int  `json:"quantity" binding:"min=1"`
	}

	if err := ctx.BindJSON(&request); err != nil {
		h.errorResponse(ctx, http.StatusBadRequest, err)
		return
	}

	flightServID, err := h.Repository.GetOrCreateCartFlightServ(userID)
	if err != nil {
		h.errorResponse(ctx, http.StatusInternalServerError, err)
		return
	}

	err = h.Repository.AddToCart(flightServID, request.ServiceID, request.Quantity)
	if err != nil {
		h.errorResponse(ctx, http.StatusInternalServerError, err)
		return
	}

	// ✅ JSON вместо Redirect
	ctx.JSON(http.StatusOK, gin.H{

		"message": "Item added to cart",
	})
}

// CreateService - создание новой услуги
func (h *Handler) CreateService(ctx *gin.Context) {
	var service ds.ASGARService

	if err := ctx.BindJSON(&service); err != nil {
		h.errorResponse(ctx, http.StatusBadRequest, err)
		return
	}

	// Устанавливаем дефолтные значения
	service.IsDelete = false

	// Сохраняем в БД
	if err := h.Repository.CreateService(&service); err != nil {
		h.errorResponse(ctx, http.StatusInternalServerError, err)
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{

		"service": service,
	})
}

// UpdateService - изменение услуги
func (h *Handler) UpdateService(ctx *gin.Context) {
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

	// ⭐ РАЗРЕШАЕМ ТОЛЬКО ОПРЕДЕЛЕННЫЕ ПОЛЯ
	allowedFields := map[string]bool{
		"name":             true,
		"info":             true,
		"price":            true,
		"full_description": true,
		"unit":             true,
	}

	// ⭐ ФИЛЬТРУЕМ - оставляем только разрешенные поля
	filteredUpdates := make(map[string]interface{})
	for key, value := range updates {
		if allowedFields[key] {
			filteredUpdates[key] = value
		}
	}

	// ⭐ ЗАПРЕЩАЕМ СИСТЕМНЫЕ ПОЛЯ
	delete(filteredUpdates, "id")
	delete(filteredUpdates, "is_delete")
	delete(filteredUpdates, "img")

	if len(filteredUpdates) == 0 {
		h.errorResponse(ctx, http.StatusBadRequest, fmt.Errorf("no valid fields to update"))
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

// DeleteService - удаление услуги (мягкое удаление)
func (h *Handler) DeleteMiniplane(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		h.errorResponse(ctx, http.StatusBadRequest, err)
		return
	}

	// Мягкое удаление - устанавливаем is_delete = true
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

// GetCartIcon - иконка корзины (id заявки + кол-во товаров)
func (h *Handler) GetMiniplaneIcon(ctx *gin.Context) {
	userID := uint(1)

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

// Вспомогательные функции
func formatPrice(price int) string {
	priceStr := strconv.Itoa(price)
	if len(priceStr) <= 3 {
		return priceStr
	}

	var result strings.Builder
	count := 0
	for i := len(priceStr) - 1; i >= 0; i-- {
		if count > 0 && count%3 == 0 {
			result.WriteString(" ")
		}
		result.WriteByte(priceStr[i])
		count++
	}

	runes := []rune(result.String())
	for i, j := 0, len(runes)-1; i < j; i, j = i+1, j-1 {
		runes[i], runes[j] = runes[j], runes[i]
	}

	return string(runes)
}

func (h *Handler) errorResponse(ctx *gin.Context, statusCode int, err error) {
	logrus.Error(err.Error())
	ctx.JSON(statusCode, gin.H{
		"status": "error",
		"error":  err.Error(),
	})
}

// GetRequests - список заявок
func (h *Handler) GetRequests(ctx *gin.Context) {
	requests, err := h.Repository.GetAllRequests()
	if err != nil {
		h.errorResponse(ctx, http.StatusInternalServerError, err)
		return
	}

	// Фильтруем удаленные и черновики
	var filteredRequests []ds.FlightServ
	for _, req := range requests {
		if req.Status != "удален" && req.Status != "черновик" {
			filteredRequests = append(filteredRequests, req)
		}
	}

	ctx.JSON(http.StatusOK, gin.H{

		"requests": filteredRequests,
		"count":    len(filteredRequests),
	})
}

// GetRequest - одна заявка
func (h *Handler) GetRequest(ctx *gin.Context) {
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

	ctx.JSON(http.StatusOK, gin.H{

		"request": gin.H{
			"id":          request.ID,
			"status":      request.Status,
			"total_price": request.TotalPrice,
			"date_create": request.DateCreate,
			"date_update": request.DateUpdate,
			"date_finish": request.DateFinish,
			"user_login":  user.Login,
			"items":       request.Subjservs,
		},
	})
}

// UpdateRequest - изменение полей заявки
func (h *Handler) UpdateRequest(ctx *gin.Context) {
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

	// Запрещаем обновление системных полей
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

// FormRequest - формирование заявки создателем
func (h *Handler) FormRequest(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		h.errorResponse(ctx, http.StatusBadRequest, err)
		return
	}

	// TODO: Проверка обязательных полей
	// TODO: Бизнес-логика смены статуса

	if err := h.Repository.UpdateRequestStatus(uint(id), "формирование"); err != nil {
		h.errorResponse(ctx, http.StatusInternalServerError, err)
		return
	}

	ctx.JSON(http.StatusOK, gin.H{

		"message": "Request formed successfully",
	})
}

// CompleteRequest - завершение заявки модератором
func (h *Handler) CompleteRequest(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		h.errorResponse(ctx, http.StatusBadRequest, err)
		return
	}

	// TODO: Вычисление стоимости по формуле из лаб 2
	// TODO: Установка модератора

	if err := h.Repository.CompleteRequest(uint(id)); err != nil {
		h.errorResponse(ctx, http.StatusInternalServerError, err)
		return
	}

	ctx.JSON(http.StatusOK, gin.H{

		"message": "Request completed",
	})
}

// DeleteRequestItem - удаление позиции из заявки (М-М)
func (h *Handler) DeleteRequestItem(ctx *gin.Context) {
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

		"message": "Item removed from request",
	})
}

// UpdateRequestItem - изменение позиции в заявке (М-М)
func (h *Handler) UpdateRequestItem(ctx *gin.Context) {
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

		"message": "Request item updated",
	})
}

// RegisterUser - регистрация пользователя
// RegisterUser - регистрация пользователя
func (h *Handler) RegisterUser(ctx *gin.Context) {
	var request struct {
		Login    string `json:"login" binding:"required,min=3"`
		Password string `json:"password" binding:"required,min=6"`
	}

	if err := ctx.BindJSON(&request); err != nil {
		h.errorResponse(ctx, http.StatusBadRequest, err)
		return
	}

	// Проверяем нет ли пользователя с таким логином
	var existingUser ds.Avius
	if err := h.Repository.GetUserByLogin(request.Login, &existingUser); err == nil {
		h.errorResponse(ctx, http.StatusBadRequest, fmt.Errorf("user with this login already exists"))
		return
	}

	// Создаем пользователя
	user := ds.Avius{
		Login:    request.Login,
		Password: request.Password, // TODO: хэшировать пароль

		IsModerator: false,
	}

	if err := h.Repository.CreateUser(&user); err != nil {
		h.errorResponse(ctx, http.StatusInternalServerError, err)
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{

		"message": "User registered successfully",
		"user_id": user.ID,
	})
}

// GetUserProfile - профиль пользователя
// GetUserProfile - профиль пользователя (с параметром id)
func (h *Handler) GetUserProfile(ctx *gin.Context) {
	idStr := ctx.Query("id")
	if idStr == "" {
		idStr = "1" // значение по умолчанию
	}

	id, err := strconv.Atoi(idStr)
	if err != nil {
		h.errorResponse(ctx, http.StatusBadRequest, err)
		return
	}

	user, err := h.Repository.GetUserByID(uint(id))
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

// UpdateUserProfile - изменение профиля (тоже с параметром id)
func (h *Handler) UpdateUserProfile(ctx *gin.Context) {
	idStr := ctx.Query("id")
	if idStr == "" {
		idStr = "1" // значение по умолчанию
	}

	id, err := strconv.Atoi(idStr)
	if err != nil {
		h.errorResponse(ctx, http.StatusBadRequest, err)
		return
	}

	// Пока просто возвращаем успех
	ctx.JSON(http.StatusOK, gin.H{

		"message": "Profile updated",
		"user_id": id,
	})
}

// LoginUser - аутентификация
func (h *Handler) LoginUser(ctx *gin.Context) {
	var request struct {
		Login    string `json:"login" binding:"required"`
		Password string `json:"password" binding:"required"`
	}

	if err := ctx.BindJSON(&request); err != nil {
		h.errorResponse(ctx, http.StatusBadRequest, err)
		return
	}

	// Проверяем пользователя
	var user ds.Avius
	if err := h.Repository.GetUserByLogin(request.Login, &user); err != nil {
		h.errorResponse(ctx, http.StatusUnauthorized, fmt.Errorf("invalid login or password"))
		return
	}

	// Проверяем пароль (пока без хэширования)
	if user.Password != request.Password {
		h.errorResponse(ctx, http.StatusUnauthorized, fmt.Errorf("invalid login or password"))
		return
	}

	// TODO: Генерация JWT токена
	ctx.JSON(http.StatusOK, gin.H{

		"message": "Login successful",
		"user_id": user.ID,
	})
}

// LogoutUser - деавторизация
func (h *Handler) LogoutUser(ctx *gin.Context) {
	// TODO: Инвалидация токена
	ctx.JSON(http.StatusOK, gin.H{

		"message": "Logout successful",
	})
}
