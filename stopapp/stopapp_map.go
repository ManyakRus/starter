package stopapp

import (
	ordered_map "github.com/m-murad/ordered-sync-map"
	"golang.org/x/net/context"
	"sync"
)

// WaitGroupContext - структура для хранения WaitGroup и контекста
type WaitGroupContext struct {
	WaitGroup     *sync.WaitGroup
	Ctx           *context.Context
	CancelCtxFunc func()
}

// OrderedMapConnections - содержит все WaitGroup от разных компонент, в порядке подключения компонентов
var OrderedMapConnections *ordered_map.Map[string, WaitGroupContext]

func init() {
	OrderedMapConnections = ordered_map.New[string, WaitGroupContext]()
}
