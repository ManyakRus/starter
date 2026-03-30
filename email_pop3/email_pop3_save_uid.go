package email_pop3

import (
	"encoding/json"
	"github.com/ManyakRus/starter/log"
	"os"
	"sync"
)

var ProcessedUIDsFilename = ""
var ProcessedUIDsFilename_short = "pop3_processed.json"

// ProcessedUIDs - хранилище обработанных UID
var ProcessedUIDs = struct {
	sync.RWMutex
	uids map[string]bool
}{
	uids: make(map[string]bool),
}

// LoadProcessedUIDs - загружает обработанные UID из файла
func LoadProcessedUIDs() error {
	data, err := os.ReadFile(ProcessedUIDsFilename)
	if err != nil {
		if os.IsNotExist(err) {
			log.Infof("No processed UIDs file, starting fresh. Filename: %s", ProcessedUIDsFilename)
			return nil
		}
		return err
	}

	var uids []string
	if err := json.Unmarshal(data, &uids); err != nil {
		return err
	}

	ProcessedUIDs.Lock()
	defer ProcessedUIDs.Unlock()
	for _, uid := range uids {
		ProcessedUIDs.uids[uid] = true
	}

	log.Infof("Loaded %d processed UIDs from %s", len(uids), ProcessedUIDsFilename)
	return nil
}

// SaveProcessedUIDs - сохраняет обработанные UID в файл
func SaveProcessedUIDs() error {
	ProcessedUIDs.RLock()
	uids := make([]string, 0, len(ProcessedUIDs.uids))
	for uid := range ProcessedUIDs.uids {
		uids = append(uids, uid)
	}
	ProcessedUIDs.RUnlock()

	data, err := json.MarshalIndent(uids, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(ProcessedUIDsFilename, data, 0644)
}

// MarkUIDAsProcessed - отмечает UID как обработанный
func MarkUIDAsProcessed(uid string) {
	if uid == "" {
		return
	}
	ProcessedUIDs.Lock()
	defer ProcessedUIDs.Unlock()
	ProcessedUIDs.uids[uid] = true
}

// IsUIDProcessed - проверяет, обработан ли UID
func IsUIDProcessed(uid string) bool {
	if uid == "" {
		return false
	}
	ProcessedUIDs.RLock()
	defer ProcessedUIDs.RUnlock()
	return ProcessedUIDs.uids[uid]
}

// CleanOldUIDs - удаляет UID (если нужно, но обычно не вызывается)
// ВНИМАНИЕ: после удаления UID письмо будет обработано повторно!
func CleanOldUIDs() {
	// По умолчанию ничего не удаляем
	// Если нужно очистить, раскомментируйте и реализуйте логику
	// но это приведёт к повторной обработке писем
}
