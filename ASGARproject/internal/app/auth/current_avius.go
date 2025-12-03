package auth

import "sync"

var (
	currentUserID uint = 1 // Константа пользователя по умолчанию
	mu            sync.RWMutex
)

// GetCurrentUserID - функция-синглтон для получения ID текущего пользователя
// Используется во всех методах, где нужен userID
func GetCurrentUserID() uint {
	mu.RLock()
	defer mu.RUnlock()
	return currentUserID
}

// SetCurrentUserIDForTesting - для тестов (может пригодиться)
func SetCurrentUserIDForTesting(id uint) {
	mu.Lock()
	defer mu.Unlock()
	currentUserID = id
}
