package repository

import (
	"crypto/sha1" // ДОБАВИТЬ
	"decode/internal/app/ds"
	"encoding/hex" // ДОБАВИТЬ
	"fmt"
	// УДАЛИТЬ: "golang.org/x/crypto/bcrypt"
)

// GetUserByLogin возвращает пользователя по логину
func (r *Repository) GetUserByLogin(login string) (*ds.Avius, error) {
	var user ds.Avius
	err := r.db.Where("login = ?", login).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// CreateUserWithRole создает пользователя с указанной ролью
func (r *Repository) CreateUserWithRole(user *ds.Avius) error {
	return r.db.Create(user).Error
}

// VerifyUser проверяет логин и пароль
func (r *Repository) VerifyUser(login, password string) (*ds.Avius, error) {
	var user ds.Avius
	err := r.db.Where("login = ?", login).First(&user).Error

	if err != nil {
		return nil, fmt.Errorf("пользователь не найден")
	}

	// ИЗМЕНИТЬ: вместо user.CheckPassword(password) используем SHA1
	h := sha1.New()
	h.Write([]byte(password))
	hashedPassword := hex.EncodeToString(h.Sum(nil))

	if user.Password != hashedPassword {
		return nil, fmt.Errorf("неверный пароль")
	}

	return &user, nil
}

// UpdateUserRole обновляет роль пользователя
func (r *Repository) UpdateUserRole(userID uint, role string) error {
	return r.db.Model(&ds.Avius{}).Where("id = ?", userID).Update("role", role).Error
}

// GetAllUsers возвращает всех пользователей
func (r *Repository) GetAllUsers() ([]ds.Avius, error) {
	var users []ds.Avius
	err := r.db.Find(&users).Error
	return users, err
}

// ИЗМЕНИТЬ HashPassword на SHA1
func HashPassword(password string) (string, error) {
	// ЗАМЕНИТЬ bcrypt на SHA1
	h := sha1.New()
	h.Write([]byte(password))
	hashedPassword := hex.EncodeToString(h.Sum(nil))
	return hashedPassword, nil
}

// ИЗМЕНИТЬ CheckPasswordHash на SHA1
func CheckPasswordHash(password, hash string) bool {
	// ЗАМЕНИТЬ bcrypt на SHA1
	h := sha1.New()
	h.Write([]byte(password))
	hashedPassword := hex.EncodeToString(h.Sum(nil))
	return hash == hashedPassword
}
