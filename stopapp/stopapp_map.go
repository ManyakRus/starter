package stopapp

import (
	ordered_map "github.com/m-murad/ordered-sync-map"
)

// MapWaitGroups - содержит все WaitGroup от разных компонент, в порядке подключения компонентов
var MapWaitGroups *ordered_map.Map[string, IWait]

func init() {
	MapWaitGroups = ordered_map.New[string, IWait]()
}
