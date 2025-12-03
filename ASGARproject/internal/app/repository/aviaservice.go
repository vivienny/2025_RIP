package repository

import (
	"decode/internal/app/ds"
)

func (r *Repository) GetAllServices() ([]ds.ASGARService, error) {
	var services []ds.ASGARService
	err := r.db.Where("is_delete = ?", false).Find(&services).Error
	if err != nil {
		return nil, err
	}
	return services, nil
}

func (r *Repository) SearchServicesByName(name string) ([]ds.ASGARService, error) {
	var services []ds.ASGARService
	query := r.db.Where("is_delete = ?", false)

	// Если name не пустой - добавляем фильтр
	if name != "" {
		query = query.Where("name ILIKE ?", "%"+name+"%")
	}

	// Всегда сортируем по ID
	err := query.Order("id ASC").Find(&services).Error
	if err != nil {
		return nil, err
	}
	return services, nil
}
