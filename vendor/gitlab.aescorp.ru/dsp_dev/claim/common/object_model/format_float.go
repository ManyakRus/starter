package object_model

import (
	"fmt"
	"strings"
)

// FormatFloat -- представление вещественных чисел в разных
type FormatFloat struct {
	val float64
}

// String -- возвращает специальное строковое представление для фронта
func (sf *FormatFloat) String() string {
	tmp := fmt.Sprintf("%d", int64(sf.val))
	res := ""
	j := 0
	for i := len(tmp) - 1; i >= 0; i-- {
		j++
		res = string(tmp[i]) + res
		if j > 0 && j%3 == 0 {
			res = " " + res
		}
	}

	tmp1 := fmt.Sprintf("%.2f", sf.val)
	tmp2 := strings.Split(tmp1, ".")
	res = res + "." + tmp2[1]

	// fmt.Println(res)
	res = strings.Trim(res, " ")
	return strings.ReplaceAll(res, ".", ",")
}

// Get -- возвращает хранимое значение
func (sf *FormatFloat) Get() float64 {
	return sf.val
}
