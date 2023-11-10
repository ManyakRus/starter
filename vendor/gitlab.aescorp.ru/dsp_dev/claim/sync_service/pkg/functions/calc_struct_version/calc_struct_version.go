package calc_struct_version

import (
	"gitlab.aescorp.ru/dsp_dev/claim/nikitin/micro"
	"reflect"
)

// CalcStructVersion - вычисляет версию модели
func CalcStructVersion(t reflect.Type) uint32 {
	var ReturnVar uint32

	names := make([]string, t.NumField())

	// имя + тип поля
	s := ""
	for i := range names {
		s = s + t.Field(i).Name
		s = s + t.Field(i).Type.Name()
	}

	ReturnVar = micro.Hash(s)

	return ReturnVar
}
