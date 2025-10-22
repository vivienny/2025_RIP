package repository

import (
	"decode/internal/app/ds"
	"fmt"

	"github.com/sirupsen/logrus"
)

func (r *Repository) GetAsgarServices() ([]ds.ASGARService, error) {
	var services []ds.ASGARService
	err := r.db.Find(&services).Error
	if err != nil {
		return nil, err
	}

	if len(services) == 0 {
		return nil, fmt.Errorf("массив пустой")
	}
	return services, nil
}

func (r *Repository) GetAsgarService(id int) (ds.ASGARService, error) {
	service := ds.ASGARService{}
	err := r.db.Where("id = ?", id).First(&service).Error
	if err != nil {
		return ds.ASGARService{}, err
	}
	return service, nil
}

func (r *Repository) GetServicesByTitle(title string) ([]ds.ASGARService, error) {
	var services []ds.ASGARService
	err := r.db.Where("title ILIKE ?", "%"+title+"%").Find(&services).Error
	if err != nil {
		return nil, err
	}
	return services, nil
}

// CreateService создает новую услугу
func (r *Repository) CreateService(service *ds.ASGARService) error {
	logrus.Printf("📝 REPOSITORY: Создание новой услуги: %s", service.Name)

	err := r.db.Create(service).Error
	if err != nil {
		logrus.Printf("❌ REPOSITORY: Ошибка создания услуги: %v", err)
		return fmt.Errorf("ошибка при создании услуги: %w", err)
	}

	logrus.Printf("✅ REPOSITORY: Услуга создана, ID=%d", service.ID)
	return nil
}

// UpdateService обновляет услугу
func (r *Repository) UpdateService(id uint, updates map[string]interface{}) error {
	logrus.Printf("📝 REPOSITORY: Обновление услуги ID=%d", id)
	logrus.Printf("📝 REPOSITORY: Updates map: %+v", updates)

	// СОХРАНЯЕМ РЕЗУЛЬТАТ В ПЕРЕМЕННУЮ
	result := r.db.Model(&ds.ASGARService{}).Where("id = ?", id).Updates(updates)

	// ДОБАВЛЯЕМ ОТЛАДОЧНУЮ ИНФОРМАЦИЮ
	logrus.Printf("📝 REPOSITORY: Rows affected: %d", result.RowsAffected)
	logrus.Printf("📝 REPOSITORY: Error: %v", result.Error)

	if result.Error != nil {
		logrus.Printf("❌ REPOSITORY: Ошибка обновления услуги ID=%d: %v", id, result.Error)
		return fmt.Errorf("ошибка при обновлении услуги: %w", result.Error)
	}

	logrus.Printf("✅ REPOSITORY: Услуга обновлена, ID=%d", id)
	return nil
}

// DeleteService удаляет услугу (мягкое удаление)
func (r *Repository) DeleteService(id uint) error {
	logrus.Printf("🗑️ REPOSITORY: Удаление услуги ID=%d", id)

	// Мягкое удаление - устанавливаем is_delete = true
	err := r.db.Model(&ds.ASGARService{}).Where("id = ?", id).Update("is_delete", true).Error
	if err != nil {
		logrus.Printf("❌ REPOSITORY: Ошибка удаления услуги ID=%d: %v", id, err)
		return fmt.Errorf("ошибка при удалении услуги: %w", err)
	}

	logrus.Printf("✅ REPOSITORY: Услуга удалена (is_delete=true), ID=%d", id)
	return nil
}
