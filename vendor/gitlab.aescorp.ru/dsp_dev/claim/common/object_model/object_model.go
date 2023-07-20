package object_model

import (
	"fmt"
	"strings"
	"time"
)

// Поиск по истории дела обновление статуса на определенный код МОЖЕТ ПРИГОДИТЬСЯ В БУДУЩЕМ
/*func isIncludedStage(briefCaseChanges []ChangeItem, desiredValue int) bool {
	for i := 0; i < len(briefCaseChanges); i++ {
		if briefCaseChanges[i].Key == "Обновление статуса" {
			value := regexp.MustCompile(`\d`).FindStringSubmatch(briefCaseChanges[i].Value)
			if len(value) == 1 {
				valueCode, _ := strconv.Atoi(value[0])
				if desiredValue == valueCode {
					return true
				}
			}
		}
	}
	return false
}*/

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

func formatDate(date time.Time) string {
	return date.Format("02.01.2006")
}

func formatTime(date time.Time) string {
	return date.Format("02.01.2006 15:04:05")
}
