package repository

import (
	"decode/internal/app/ds"
	"fmt"
	"strconv"

	"github.com/sirupsen/logrus"
)

// GetMiniplaneItemsCount для получения количества услуг в корзине
func (r *Repository) GetMiniplaneItemsCount() int64 {
	var miniplaneID uint
	var count int64
	userID := uint(1) // ← ИСПРАВЛЕНО НА uint(1)

	logrus.Printf("🔢 REPOSITORY: Поиск корзины для подсчета, UserID=%d", userID)

	// Ищем активную корзину пользователя
	err := r.db.Model(&ds.Miniplane{}).Where("user_id = ? AND is_active = ?", userID, true).Select("id").First(&miniplaneID).Error
	if err != nil {
		logrus.Printf("❌ REPOSITORY: Корзина не найдена для подсчета, UserID=%d", userID)
		return 0
	}

	logrus.Printf("🔢 REPOSITORY: Корзина найдена для подсчета, MiniplaneID=%d", miniplaneID)

	// Считаем позиции в корзине
	err = r.db.Model(&ds.Subjserv{}).Where("miniplane_id = ?", miniplaneID).Count(&count).Error
	if err != nil {
		logrus.Printf("❌ REPOSITORY: Ошибка подсчета позиций: %v", err)
		return 0
	}

	logrus.Printf("🔢 REPOSITORY: В корзине %d позиций", count)
	return count
}

// GetActiveMiniplaneID получает ID активной корзины пользователя
func (r *Repository) GetActiveMiniplaneID(userID uint) (uint, error) {
	var miniplane ds.Miniplane
	logrus.Printf("🔍 REPOSITORY: Поиск активной корзины для UserID=%d", userID)

	err := r.db.Where("user_id = ? AND is_active = ?", userID, true).First(&miniplane).Error
	if err != nil {
		logrus.Printf("❌ REPOSITORY: Активная корзина не найдена для UserID=%d, ошибка: %v", userID, err)
		return 0, err
	}

	logrus.Printf("✅ REPOSITORY: Активная корзина найдена, ID=%d", miniplane.ID)
	return uint(miniplane.ID), nil
}

// GetMiniplaneItems получает все позиции корзины
func (r *Repository) GetMiniplaneItems(miniplaneID uint) ([]ds.Subjserv, error) {
	var items []ds.Subjserv
	logrus.Printf("📦 REPOSITORY: Получение позиций корзины MiniplaneID=%d", miniplaneID)

	err := r.db.Where("miniplane_id = ?", miniplaneID).Preload("Service").Find(&items).Error
	if err != nil {
		logrus.Printf("❌ REPOSITORY: Ошибка получения позиций: %v", err)
		return nil, err
	}

	logrus.Printf("✅ REPOSITORY: Получено %d позиций для MiniplaneID=%d", len(items), miniplaneID)
	for i, item := range items {
		logrus.Printf("   🛒 Позиция %d: ServiceID=%d, Service.Name=%s", i+1, item.ServiceID, item.Service.Name)
	}

	return items, nil
}

// GetMiniplaneTotal вычисляет общую сумму корзины на лету
func (r *Repository) GetMiniplaneTotal(miniplaneID uint) (int, error) {
	var items []ds.Subjserv
	logrus.Printf("💰 REPOSITORY: Вычисление суммы корзины MiniplaneID=%d", miniplaneID)

	// Получаем все позиции корзины с услугами
	err := r.db.Where("miniplane_id = ?", miniplaneID).Preload("Service").Find(&items).Error
	if err != nil {
		logrus.Printf("❌ REPOSITORY: Ошибка получения позиций корзины: %v", err)
		return 0, err
	}

	// Считаем сумму в коде
	total := 0
	for _, item := range items {
		// Преобразуем строку цены в число
		price, err := strconv.Atoi(item.Service.Price)
		if err != nil {
			logrus.Printf("⚠️ REPOSITORY: Ошибка преобразования цены '%s': %v", item.Service.Price, err)
			continue
		}
		total += item.Quantity * price
	}

	logrus.Printf("💰 REPOSITORY: Сумма корзины MiniplaneID=%d: %d", miniplaneID, total)
	return total, nil
}

// RemoveFromMiniplane удаляет позицию из корзины
func (r *Repository) RemoveFromMiniplane(subjservID uint) error {
	logrus.Printf("🗑️ REPOSITORY: Удаление позиции из корзины, SubjservID=%d", subjservID)

	err := r.db.Where("id = ?", subjservID).Delete(&ds.Subjserv{}).Error
	if err != nil {
		logrus.Printf("❌ REPOSITORY: Ошибка удаления позиции: %v", err)
		return fmt.Errorf("ошибка при удалении из корзины с id %d: %w", subjservID, err)
	}

	logrus.Printf("✅ REPOSITORY: Позиция удалена, SubjservID=%d", subjservID)
	return nil
}

// ClearMiniplane очищает всю корзину
func (r *Repository) ClearMiniplane(miniplaneID uint) error {
	logrus.Printf("🗑️ REPOSITORY: Очистка всей корзины, MiniplaneID=%d", miniplaneID)

	err := r.db.Where("miniplane_id = ?", miniplaneID).Delete(&ds.Subjserv{}).Error
	if err != nil {
		logrus.Printf("❌ REPOSITORY: Ошибка очистки корзины: %v", err)
		return fmt.Errorf("ошибка при очистке корзины: %w", err)
	}

	logrus.Printf("✅ REPOSITORY: Корзина очищена, MiniplaneID=%d", miniplaneID)
	return nil
}

// UpdateSubjservQuantity обновляет количество в позиции корзины
func (r *Repository) UpdateSubjservQuantity(subjservID uint, quantity int) error {
	logrus.Printf("📝 REPOSITORY: Обновление количества SubjservID=%d, Quantity=%d", subjservID, quantity)

	err := r.db.Model(&ds.Subjserv{}).Where("id = ?", subjservID).Update("quantity", quantity).Error
	if err != nil {
		logrus.Printf("❌ REPOSITORY: Ошибка обновления количества: %v", err)
		return fmt.Errorf("ошибка при обновлении количества: %w", err)
	}

	logrus.Printf("✅ REPOSITORY: Количество обновлено, SubjservID=%d", subjservID)
	return nil
}
