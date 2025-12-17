package stopapp

import "sync"

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
