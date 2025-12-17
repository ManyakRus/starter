package postgres_sqlx

import (
	"golang.org/x/net/context"
	"sync"
)

// ctx_Connect, cancelCtxFunc - контекст для одного соединения, при отмене контекста соединение закроется
var ctx_Connect, cancelCtxFunc = context.WithCancel(context.Background())

// waitGroup_Connect - группа ожидания завершения всех частей программы
var waitGroup_Connect *sync.WaitGroup

// lockWaitGroup_Connect - гарантирует получение WGMain с учётом многопоточности
var lockWaitGroup_Connect = &sync.RWMutex{}

// SetWaitGroup_Connect - присваивает внешний WaitGroup
func SetWaitGroup_Connect(wg *sync.WaitGroup) {
	lockWaitGroup_Connect.RLock()
	defer lockWaitGroup_Connect.RUnlock()

	waitGroup_Connect = wg
}

// GetWaitGroup_Connect - возвращает группу ожидания завершения всех частей программы
func GetWaitGroup_Connect() *sync.WaitGroup {
	lockWaitGroup_Connect.Lock()
	defer lockWaitGroup_Connect.Unlock()

	if waitGroup_Connect == nil {
		waitGroup_Connect = &sync.WaitGroup{}
	}

	return waitGroup_Connect
}
