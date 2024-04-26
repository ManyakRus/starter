package calc_struct_version

import (
	"github.com/ManyakRus/starter/micro"
	"reflect"
)

// CalcStructVersion - вычисляет версию модели
func CalcStructVersion(t reflect.Type) uint32 {
	var ReturnVar uint32

	names := make([]string, t.NumField())

	// имя + тип поля
	s := ""
	for i := range names {
		Field1 := t.Field(i)
		s = s + Field1.Name
		s = s + Field1.Type.Name()
		if Field1.Anonymous == true && Field1.Type != t {
			version2 := CalcStructVersion(Field1.Type)
			ReturnVar = ReturnVar + version2
		}
	}

	ReturnVar = ReturnVar + micro.Hash(s)

	return ReturnVar
}
