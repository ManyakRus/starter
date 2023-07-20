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

// lockContextMain - гарантирует единственное создание контеста
// var lockContextMain sync.Mutex

// GetContext - возвращает глобальный контекст приложения
func GetContext() context.Context {
	//lockContextMain.Lock()
	//defer lockContextMain.Unlock()
	//
	//if Ctx == nil {
	//	CtxBg := context.Background()
	//	Ctx, CancelContext = context.WithCancel(CtxBg)
	//}

	onceCtx.Do(func() {
		CtxBg := context.Background()
		var Ctx0 context.Context
		Ctx0, CancelContext = context.WithCancel(CtxBg)
		Ctx = &Ctx0
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
