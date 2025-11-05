package handler

import (
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
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		logrus.Error(err)
		return
	}

	ctx.HTML(http.StatusOK, "asgar_catalog.html", gin.H{
		"services":   services,
		"cart_count": h.Repository.GetCartItemsCount(), // ← ИЗМЕНИЛИ
		"search":     search,
	})
}

// GetMiniplane отображает страницу корзины
func (h *Handler) GetMiniplane(ctx *gin.Context) {
	userID := uint(1) // пока захардкодили, потом из JWT

	logrus.Printf("🛒 HANDLER: Начинаем поиск корзины для UserID=%d", userID)

	// Получаем активную корзину-заявку
	flightServID, err := h.Repository.GetOrCreateCartFlightServ(userID) // ← ИЗМЕНИЛИ
	if err != nil {
		logrus.Printf("❌ HANDLER: Корзина не найдена для UserID=%d, ошибка: %v", userID)
		// Если корзины нет - создаем пустую
		ctx.HTML(http.StatusOK, "miniplane.html", gin.H{
			"items": []ds.Subjserv{},
			"total": "0",
		})
		return
	}

	logrus.Printf("✅ HANDLER: Корзина найдена, FlightServID=%d", flightServID)

	// Получаем позиции корзины
	items, err := h.Repository.GetCartItems(flightServID) // ← ИЗМЕНИЛИ
	if err != nil {
		logrus.Printf("❌ HANDLER: Ошибка получения позиций корзины: %v", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		logrus.Error(err)
		return
	}

	logrus.Printf("✅ HANDLER: Получено позиций в корзине: %d", len(items))
	for i, item := range items {
		logrus.Printf("   📦 Позиция %d: ServiceID=%d, Quantity=%d, Price=%s",
			i+1, item.ServiceID, item.Quantity, item.Price)
	}

	// Считаем общую сумму
	total := 0
	for _, item := range items {
		// Убираем пробелы и преобразуем цену в число
		cleanPrice := strings.ReplaceAll(item.Price, " ", "")
		price, err := strconv.Atoi(cleanPrice)
		if err != nil {
			logrus.Printf("⚠️ Ошибка преобразования цены: %s", item.Price)
			continue
		}
		total += price * item.Quantity
	}
	totalStr := formatPrice(total)

	logrus.Printf("💰 Общая стоимость корзины: %d руб.", total)

	ctx.HTML(http.StatusOK, "miniplane.html", gin.H{
		"items": items,
		"total": totalStr,
	})
}

// RemoveFromMiniplane удаляет позицию из корзины
func (h *Handler) RemoveFromMiniplane(ctx *gin.Context) {
	strId := ctx.PostForm("subjserv_id")
	id, err := strconv.Atoi(strId)
	if err != nil {
		logrus.Printf("❌ HANDLER: Неверный ID позиции: %s", strId)
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	logrus.Printf("🗑️ HANDLER: Удаление позиции из корзины, SubjservID=%d", id)

	err = h.Repository.RemoveFromCart(uint(id)) // ← ИЗМЕНИЛИ
	if err != nil {
		logrus.Printf("❌ HANDLER: Ошибка удаления позиции: %v", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	logrus.Printf("✅ HANDLER: Позиция удалена, перенаправление в корзину")
	// Перенаправляем обратно в корзину
	ctx.Redirect(http.StatusFound, "/miniplane")
}

// ClearMiniplane очищает всю корзину
func (h *Handler) ClearMiniplane(ctx *gin.Context) {
	userID := uint(1)

	logrus.Printf("🗑️ HANDLER: Очистка корзины для UserID=%d", userID)

	// Получаем активную корзину-заявку
	flightServID, err := h.Repository.GetOrCreateCartFlightServ(userID) // ← ИЗМЕНИЛИ
	if err != nil {
		logrus.Printf("❌ HANDLER: Корзина не найдена для очистки")
		ctx.Redirect(http.StatusFound, "/miniplane")
		return
	}

	err = h.Repository.ClearCart(flightServID) // ← ИЗМЕНИЛИ
	if err != nil {
		logrus.Printf("❌ HANDLER: Ошибка очистки корзины: %v", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	logrus.Printf("✅ HANDLER: Корзина очищена, перенаправление")
	ctx.Redirect(http.StatusFound, "/miniplane")
}

// GetAsgarService получает конкретную услугу
func (h *Handler) GetAsgarService(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	service, err := h.Repository.GetAsgarService(id)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "Service not found"})
		return
	}

	ctx.HTML(http.StatusOK, "selected_aviaproduct.html", gin.H{
		"service": service,
	})
}

// formatPrice форматирует цену с пробелами (оставляем без изменений)
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

	// Переворачиваем обратно
	runes := []rune(result.String())
	for i, j := 0, len(runes)-1; i < j; i, j = i+1, j-1 {
		runes[i], runes[j] = runes[j], runes[i]
	}

	return string(runes)
}
