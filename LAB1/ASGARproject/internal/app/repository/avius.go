package repository

import (
	"decode/internal/app/ds"
)

// AviusRepository репозиторий для работы с пользователями
func (r *Repository) CreateUser(user *ds.Avius) error {
	return r.db.Create(user).Error
}

func (r *Repository) GetUserByLogin(login string) (ds.Avius, error) {
	var user ds.Avius
	err := r.db.Where("login = ?", login).First(&user).Error
	return user, err
}

func (r *Repository) GetUserByID(id uint) (ds.Avius, error) {
	var user ds.Avius
	err := r.db.First(&user, id).Error
	return user, err
}

func (r *Repository) UpdateUser(user *ds.Avius) error {
	return r.db.Save(user).Error
}
