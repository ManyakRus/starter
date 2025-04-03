// модуль для получения единого Context приложения

package contextmain

import (
	"context"
	"sync"
)

// Ctx хранит глобальный контекст программы
// не использовать
// * - чтоб можно было засунуть ссылку на чужой контекст
var Ctx *context.Context

// CancelContext - функция отмены глобального контекста
var CancelContext func()

// onceCtx - гарантирует единственное создание контеста
var onceCtx sync.Once

// MutexContextMain - гарантирует единственное создание контеста
var MutexContextMain sync.RWMutex

// GetContext - возвращает глобальный контекст приложения
func GetContext() context.Context {
	MutexContextMain.RLock()
	defer MutexContextMain.RUnlock()

	//if Ctx == nil {
	//	CtxBg := context.Background()
	//	Ctx, CancelContext = context.WithCancel(CtxBg)
	//}

	onceCtx.Do(func() {
		if Ctx == nil { //можно заполнить свой контекст, поэтому if
			CtxBg := context.Background()
			var Ctx0 context.Context
			Ctx0, CancelContext = context.WithCancel(CtxBg)
			Ctx = &Ctx0
		}
	})

	return *Ctx
}

// GetNewContext - создаёт и возвращает новый контекст приложения
func GetNewContext() context.Context {
	CtxBg := context.Background()
	var Ctx0 context.Context
	Ctx0, CancelContext = context.WithCancel(CtxBg)
	Ctx = &Ctx0

	return *Ctx
}

// SetContext - устанавливает глобальный контекст, с учётом Mutex
func SetContext(ctx *context.Context) {
	MutexContextMain.Lock()
	defer MutexContextMain.Unlock()
	Ctx = ctx
}

// SetCancelContext - устанавливает функцию глобального отмены контекста, с учётом Mutex
func SetCancelContext(cancelContext func()) {
	MutexContextMain.Lock()
	defer MutexContextMain.Unlock()
	CancelContext = cancelContext
}
