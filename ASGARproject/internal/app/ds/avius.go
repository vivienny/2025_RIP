package ds

import (
	"crypto/sha1"
	"encoding/hex"
	"strings"

	"gorm.io/gorm"
)

type Avius struct {
	ID          uint   `gorm:"primaryKey"`
	Login       string `gorm:"type:varchar(25);unique;not null"`
	Password    string `gorm:"type:varchar(100);not null"`
	IsModerator bool   `gorm:"type:boolean;default:false"`
	Role        string `gorm:"type:varchar(20);default:'user'"`
}

// Проверка пароля с SHA-1
func (u *Avius) CheckPassword(password string) bool {
	// Очищаем пароль от пробелов
	password = strings.TrimSpace(password)

	// Создаем SHA-1 хеш
	hash := sha1.Sum([]byte(password))

	// Преобразуем в hex строку
	hashedPassword := hex.EncodeToString(hash[:])

	// Сравниваем с сохраненным хешем
	return strings.ToLower(u.Password) == strings.ToLower(hashedPassword)
}

// Хеширование пароля перед созданием
func (u *Avius) BeforeCreate(tx *gorm.DB) error {
	if u.Role == "" {
		u.Role = "user"
	}

	// Хешируем пароль при создании
	if u.Password != "" {
		// Очищаем от пробелов
		u.Password = strings.TrimSpace(u.Password)

		// Проверяем, не хеширован ли уже пароль
		// SHA-1 хеш в hex имеет длину 40 символов
		if len(u.Password) == 40 {
			// Пробуем декодировать как hex
			_, err := hex.DecodeString(u.Password)
			if err == nil {
				// Уже валидный hex хеш, оставляем как есть
				return nil
			}
		}

		// Хешируем пароль
		hash := sha1.Sum([]byte(u.Password))
		u.Password = hex.EncodeToString(hash[:])
	}
	return nil
}
