package repository

import (
	"decode/internal/app/ds"
	"fmt"
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
