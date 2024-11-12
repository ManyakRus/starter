package sync_global

import (
	"fmt"
	"strings"
	"sync"
)

const (
	SyncRoot    = "/claim_reqrepl/"
	SyncQueue   = "sync_exchange"
	SyncDestVer = "1.0"
)

var (
	block       sync.RWMutex
	syncService = "unknown_service"
)

// SyncService -- возвращает имя сервиса для сетевого обмена по локальной шине
func SyncService() string {
	block.RLock()
	defer block.RUnlock()
	return syncService
}

// SetSyncService -- устанавливает имя сервиса для сетевого обмена по локальной шине
func SetSyncService(name string) error {
	block.Lock()
	defer block.Unlock()
	if strings.Trim(name, " ") == "" {
		return fmt.Errorf("SetSyncService(): name(%q) is bad", name)
	}
	syncService = name
	return nil
}
