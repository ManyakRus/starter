// package bus_conn -- подключение к локальной шине
package bus_conn

import (
	"sync"
)

// BusConn -- подключение к локальной шине
type BusConn struct {
	isConn bool // Признак подключения шины
	block  sync.RWMutex
}

// Set -- устанавливает состояние подключения
func (sf *BusConn) Set() {
	sf.block.Lock()
	defer sf.block.Unlock()
	sf.isConn = true
}

// Reset -- сбрасывает состояние подключения
func (sf *BusConn) Reset() {
	sf.block.Lock()
	defer sf.block.Unlock()
	sf.isConn = false
}

// IsConnect -- возвращает признак подключения
func (sf *BusConn) IsConnect() bool {
	sf.block.RLock()
	defer sf.block.RUnlock()
	return sf.isConn
}
