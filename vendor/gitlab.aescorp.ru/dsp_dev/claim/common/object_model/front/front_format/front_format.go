// package front_format -- форматирование данных для фронта
package front_format

import (
	"fmt"
	"strings"
	"time"

	"gitlab.aescorp.ru/dsp_dev/claim/common/object_model/alias"
)

// FrontDate -- форматирует дату в виде "02.01.2006" для фронта
func FrontDate(date time.Time) alias.FrontDate {
	strDate := date.Format("02.01.2006")
	return alias.FrontDate(strDate)
}

// FrontTime -- форматирует время в виде "02.01.2006 15:04:05" для фронта
func FrontTime(date time.Time) alias.FrontTime {
	strTime := date.Format("02.01.2006 15:04:05")
	return alias.FrontTime(strTime)
}

// Currency -- форматирует деньги для отображения
func Currency(number float64) string {
	tmp := fmt.Sprintf("%d", int64(number))
	res := ""
	j := 0
	for i := len(tmp) - 1; i >= 0; i-- {
		j++
		res = string(tmp[i]) + res
		if j > 0 && j%3 == 0 {
			res = " " + res
		}
	}

	tmp1 := fmt.Sprintf("%.2f", number)
	tmp2 := strings.Split(tmp1, ".")
	res = res + "." + tmp2[1]

	// fmt.Println(res)
	res = strings.Trim(res, " ")
	return strings.ReplaceAll(res, ".", ",")
}
