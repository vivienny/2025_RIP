package repository

import (
	"decode/internal/app/ds"
	"fmt"
	"time"

	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

// GetOrCreateCartFlightServ получает или создает корзину-заявку
func (r *Repository) GetOrCreateCartFlightServ(userID uint) (uint, error) {
	var flightServ ds.FlightServ

	err := r.db.Where("user_id = ? AND status = ?", userID, "черновик").First(&flightServ).Error
	if err == gorm.ErrRecordNotFound {
		flightServ = ds.FlightServ{
			UserID:     userID,
			Status:     "черновик",
			TotalPrice: "0",
			DateCreate: time.Now(),
			DateUpdate: time.Now(),
		}
		err = r.db.Create(&flightServ).Error
		if err != nil {
			return 0, fmt.Errorf("ошибка создания корзины: %w", err)
		}
		logrus.Printf("✅ Создана новая корзина-заявка, ID=%d", flightServ.ID)
	} else if err != nil {
		return 0, fmt.Errorf("ошибка поиска корзины: %w", err)
	}

	return flightServ.ID, nil
}

// GetCartItemsCount аналог GetMiniplaneItemsCount
func (r *Repository) GetCartItemsCount() int64 {
	var count int64
	userID := uint(1) // временно

	flightServID, err := r.GetOrCreateCartFlightServ(userID)
	if err != nil {
		return 0
	}

	r.db.Model(&ds.Subjserv{}).Where("flight_serv_id = ?", flightServID).Count(&count)
	return count
}

// GetFlightServByID - получить заявку по ID
func (r *Repository) GetFlightServByID(id uint) (ds.FlightServ, error) {
	var flightServ ds.FlightServ
	err := r.db.Preload("User").First(&flightServ, id).Error
	return flightServ, err
}
