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
	err := r.db.Where("name ILIKE ? AND is_delete = ?", "%"+name+"%", false).Find(&services).Error
	if err != nil {
		return nil, err
	}
	return services, nil
}
