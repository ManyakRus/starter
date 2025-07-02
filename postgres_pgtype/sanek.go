package postgres_pgtype

import "reflect"

// getInterfaceName - возвращает имя типа интерфейса
func getInterfaceName(v interface{}) string {
	return reflect.TypeOf(v).String()
}
