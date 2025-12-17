package repository

import (
	"decode/internal/app/ds"
	"fmt"
	"time"
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
func (r *Repository) UpdateService(id uint, updates map[string]interface{}) error {
	return r.db.Model(&ds.ASGARService{}).Where("id = ?", id).Updates(updates).Error
}

// В любой файл repository/ добави:
func (r *Repository) CreateService(service *ds.ASGARService) error {
	return r.db.Create(service).Error
}

// GetAllRequests - все заявки
func (r *Repository) GetAllRequests() ([]ds.FlightServ, error) {
	var requests []ds.FlightServ
	err := r.db.Preload("User").Find(&requests).Error
	return requests, err
}

// GetRequestWithItems - заявка с услугами
func (r *Repository) GetRequestWithItems(id uint) (ds.FlightServ, error) {
	var request ds.FlightServ
	err := r.db.Preload("Subjservs.Service").Preload("User").First(&request, id).Error
	return request, err
}

// GetUserByID - пользователь по ID
func (r *Repository) GetUserByID(id uint) (ds.Avius, error) {
	var user ds.Avius
	err := r.db.First(&user, id).Error
	return user, err
}

// UpdateRequest - обновление заявки
func (r *Repository) UpdateRequest(id uint, updates map[string]interface{}) error {
	return r.db.Model(&ds.FlightServ{}).Where("id = ?", id).Updates(updates).Error
}

// UpdateRequestStatus - обновление статуса заявки
func (r *Repository) UpdateRequestStatus(id uint, status string) error {
	return r.db.Model(&ds.FlightServ{}).Where("id = ?", id).Update("status", status).Error
}

// CompleteRequest - завершение заявки
func (r *Repository) CompleteRequest(id uint) error {
	updates := map[string]interface{}{
		"status":      "завершено",
		"date_finish": time.Now(),
		// TODO: установить moderator_id когда будет авторизация
	}
	return r.db.Model(&ds.FlightServ{}).Where("id = ?", id).Updates(updates).Error
}

// DeleteRequestItem - удаление позиции из заявки
func (r *Repository) DeleteRequestItem(id uint) error {
	return r.db.Delete(&ds.Subjserv{}, id).Error
}

// UpdateRequestItem - изменение количества в позиции заявки
func (r *Repository) UpdateRequestItem(id uint, quantity int) error {
	return r.db.Model(&ds.Subjserv{}).Where("id = ?", id).Update("quantity", quantity).Error
}

/* GetUserByLogin - пользователь по логину
func (r *Repository) GetUserByLogin(login string, user *ds.Avius) error {
	return r.db.Where("login = ?", login).First(user).Error
}
*/
// CreateUser - создание пользователя
func (r *Repository) CreateUser(user *ds.Avius) error {
	return r.db.Create(user).Error
}

// UpdateUser - обновление пользователя
func (r *Repository) UpdateUser(id uint, updates map[string]interface{}) error {
	return r.db.Model(&ds.Avius{}).Where("id = ?", id).Updates(updates).Error
}

// GetFilteredServices - услуги с фильтрацией (для будущего расширения)
func (r *Repository) GetFilteredServices(search string) ([]ds.ASGARService, error) {
	var services []ds.ASGARService
	query := r.db.Where("is_delete = ?", false)

	if search != "" {
		query = query.Where("name ILIKE ? OR info ILIKE ?",
			"%"+search+"%", "%"+search+"%")
	}

	err := query.Find(&services).Error
	return services, err
}
