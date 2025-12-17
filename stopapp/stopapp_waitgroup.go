package stopapp

import (
	"golang.org/x/net/context"
	"sync"
)

// ctx_Connect, cancelCtxFunc - контекст для одного соединения, при отмене контекста соединение закроется
var ctx_Connect, cancelCtxFunc = context.WithCancel(context.Background())

// wgMain - группа ожидания завершения всех частей программы (кроме подключений к внешним сервисам)
var wgMain *sync.WaitGroup

// lockWGMain - гарантирует получение WGMain с учётом многопоточности
var lockWGMain = &sync.RWMutex{}

// SetWaitGroup_Main - присваивает внешний WaitGroup
func SetWaitGroup_Main(wg *sync.WaitGroup) {
	lockWGMain.RLock()
	defer lockWGMain.RUnlock()

	wgMain = wg
}

// GetWaitGroup_Main - возвращает группу ожидания завершения всех частей программы
func GetWaitGroup_Main() *sync.WaitGroup {
	lockWGMain.Lock()
	defer lockWGMain.Unlock()

	if wgMain == nil {
		wgMain = &sync.WaitGroup{}
	}

	return wgMain
}
