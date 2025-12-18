package camunda_connect

import (
	"golang.org/x/net/context"
	"sync"
)

// waitGroup_Connect - группа ожидания завершения всех частей программы
var waitGroup_Connect = new(sync.WaitGroup)

// ctx_Connect - контекст для одного соединения, при отмене контекста соединение закроется
var ctx_Connect *context.Context

// cancelCtxFunc - функция для отмены контекста
var cancelCtxFunc func()

// mutex_WaitGroup_Connect - мьютекс для WaitGroup_Connect
var mutex_WaitGroup_Connect = sync.RWMutex{}

// mutex_Ctx_Connect - мьютекс для Ctx_Connect
var mutex_Ctx_Connect = sync.RWMutex{}

// init - инициализация переменных
func init() {
	ctx1, CancelFunc1 := context.WithCancel(context.Background())
	ctx_Connect = &ctx1
	cancelCtxFunc = CancelFunc1
}

// SetWaitGroup - присваивает внешний WaitGroup
func SetWaitGroup(wg *sync.WaitGroup) {
	mutex_WaitGroup_Connect.RLock()
	defer mutex_WaitGroup_Connect.RUnlock()

	waitGroup_Connect = wg
}

// GetWaitGroup - возвращает группу ожидания завершения всех частей программы
func GetWaitGroup() *sync.WaitGroup {
	mutex_WaitGroup_Connect.Lock()
	defer mutex_WaitGroup_Connect.Unlock()

	if waitGroup_Connect == nil {
		waitGroup_Connect = &sync.WaitGroup{}
	}

	return waitGroup_Connect
}

// GetContext возвращает указатель на контекст с защитой RLock
func GetContext() *context.Context {
	mutex_Ctx_Connect.RLock()
	defer mutex_Ctx_Connect.RUnlock()

	return ctx_Connect
}

// SetContext устанавливает новое значение контекста с защитой Lock
func SetContext(ctx *context.Context) {
	mutex_Ctx_Connect.Lock()
	defer mutex_Ctx_Connect.Unlock()

	ctx_Connect = ctx
}
