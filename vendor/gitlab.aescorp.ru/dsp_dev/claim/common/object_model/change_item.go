package object_model

import (
// "regexp"
// "strconv"
)

// ChangeItem Изменения
// Key - изменённое поле
// Value - новое значение
// Prev - прежнее значение
type ChangeItem struct {
	CommonStruct
	ExtLinkStruct
	// TODO UserID
	// TODO Action
	// TODO Table
	// TODO Field
	Key   string `json:"key"   gorm:"column:key;default:\"\""`
	Value string `json:"value" gorm:"column:value;default:\"\""`
	Prev  string `json:"prev"  gorm:"column:prev;default:\"\""`
}

// FIXME: не используется. Поиск по истории дела обновление статуса на определенный код
// func isIncludedStage(briefCaseChanges []ChangeItem, desiredValue int) bool {
// 	for i := 0; i < len(briefCaseChanges); i++ {
// 		if briefCaseChanges[i].Key == "Обновление статуса" {
// 			value := regexp.MustCompile(`\d`).FindStringSubmatch(briefCaseChanges[i].Value)
// 			if len(value) == 1 {
// 				valueCode, _ := strconv.Atoi(value[0])
// 				if desiredValue == valueCode {
// 					return true
// 				}
// 			}
// 		}
// 	}
// 	return false
// }
