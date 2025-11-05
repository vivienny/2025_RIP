package repository

import (
	"decode/internal/app/ds"
	"fmt"

	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

// calculateTotal вычисляет общую сумму заявки
func (r *Repository) calculateTotal(flightServID uint) string {
	var total float64

	// Суммируем: quantity * price для всех позиций заявки
	err := r.db.Model(&ds.Subjserv{}).
		Select("COALESCE(SUM(quantity * CAST(price AS NUMERIC)), 0)").
		Where("flight_serv_id = ?", flightServID).
		Scan(&total).Error

	if err != nil {
		logrus.Printf("❌ Ошибка расчёта суммы: %v", err)
		return "0"
	}

	return fmt.Sprintf("%.0f", total)
}

// GetCartItems получает все позиции корзины
func (r *Repository) GetCartItems(flightServID uint) ([]ds.Subjserv, error) {
	var items []ds.Subjserv
	logrus.Printf("📦 REPOSITORY: Получение позиций корзины FlightServID=%d", flightServID)

	err := r.db.Where("flight_serv_id = ?", flightServID).Preload("Service").Find(&items).Error
	if err != nil {
		logrus.Printf("❌ REPOSITORY: Ошибка получения позиций: %v", err)
		return nil, err
	}

	logrus.Printf("✅ REPOSITORY: Получено %d позиций для FlightServID=%d", len(items), flightServID)
	return items, nil
}

// AddToCart добавляет услугу в корзину с АВТОМАТИЧЕСКОЙ ценой
func (r *Repository) AddToCart(flightServID uint, serviceID uint, quantity int) error {
	logrus.Printf("🎯 AddToCart START: flight_serv_id=%d, service_id=%d, quantity=%d",
		flightServID, serviceID, quantity)

	// 1. СНАЧАЛА проверяем, есть ли уже такая услуга в корзине
	var existingItem ds.Subjserv
	err := r.db.Where("flight_serv_id = ? AND service_id = ?", flightServID, serviceID).First(&existingItem).Error

	if err == nil {
		// 2. Если УЖЕ ЕСТЬ - увеличиваем количество
		result := r.db.Model(&existingItem).Update("quantity", gorm.Expr("quantity + ?", quantity))
		if result.Error != nil {
			return fmt.Errorf("ошибка обновления количества: %w", result.Error)
		}
		logrus.Printf("🔄 Обновлено количество: +%d", quantity)
	} else {
		// 3. Если НЕТ - создаём новую запись
		var service ds.ASGARService
		if err := r.db.First(&service, serviceID).Error; err != nil {
			return fmt.Errorf("услуга не найдена: %w", err)
		}

		cartItem := ds.Subjserv{
			FlightServID: flightServID,
			ServiceID:    serviceID,
			Quantity:     quantity,
			Price:        service.Price,
		}

		if err := r.db.Create(&cartItem).Error; err != nil {
			return fmt.Errorf("ошибка добавления в корзину: %w", err)
		}
		logrus.Printf("✅ Добавлено в корзину: %s x%d", service.Name, quantity)
	}

	// 4. ✅ ПЕРЕСЧИТЫВАЕМ ОБЩУЮ СУММУ
	total := r.calculateTotal(flightServID)
	if err := r.db.Model(&ds.FlightServ{}).
		Where("id = ?", flightServID).
		Update("total_price", total).Error; err != nil {
		logrus.Printf("⚠️ Не удалось обновить total_price: %v", err)
	}

	logrus.Printf("💰 Обновлена общая сумма заявки %d: %s", flightServID, total)
	return nil
}

// RemoveFromCart удаляет позицию из корзины
func (r *Repository) RemoveFromCart(subjservID uint) error {
	logrus.Printf("🗑️ REPOSITORY: Удаление позиции из корзины, SubjservID=%d", subjservID)

	// Сначала узнаем flight_serv_id чтобы потом пересчитать total
	var item ds.Subjserv
	if err := r.db.First(&item, subjservID).Error; err != nil {
		return err
	}
	flightServID := item.FlightServID

	// Удаляем позицию
	if err := r.db.Where("id = ?", subjservID).Delete(&ds.Subjserv{}).Error; err != nil {
		return err
	}

	// ✅ ПЕРЕСЧИТЫВАЕМ СУММУ
	total := r.calculateTotal(flightServID)
	r.db.Model(&ds.FlightServ{}).Where("id = ?", flightServID).Update("total_price", total)

	logrus.Printf("💰 Обновлена сумма после удаления: %s", total)
	return nil
}

// ClearCart очищает всю корзину
func (r *Repository) ClearCart(flightServID uint) error {
	logrus.Printf("🗑️ REPOSITORY: Очистка всей корзины, FlightServID=%d", flightServID)

	// Очищаем корзину
	if err := r.db.Where("flight_serv_id = ?", flightServID).Delete(&ds.Subjserv{}).Error; err != nil {
		return err
	}

	// ✅ ОБНУЛЯЕМ СУММУ
	r.db.Model(&ds.FlightServ{}).Where("id = ?", flightServID).Update("total_price", "0")

	logrus.Printf("💰 Корзина очищена, сумма обнулена")
	return nil
}
