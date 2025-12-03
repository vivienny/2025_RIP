package handler

import "sync"

var (
	currentUserID uint = 1
	mu            sync.RWMutex
)

func GetCurrentUserID() uint {
	mu.RLock()
	defer mu.RUnlock()
	return currentUserID
}
