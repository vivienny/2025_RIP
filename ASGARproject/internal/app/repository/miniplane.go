package repository

import (
	"decode/internal/app/ds"
	"fmt"

	"github.com/sirupsen/logrus"
)

// GetMiniplaneItemsCount для получения количества услуг в корзине
func (r *Repository) GetMiniplaneItemsCount() int64 {
	var miniplaneID uint
	var count int64
	userID := 1 // пока захардкодили, потом из JWT

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
