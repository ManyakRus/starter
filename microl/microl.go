// пакет для микрофункций с логгером

package microl

import (
	"fmt"
	"github.com/ManyakRus/starter/log"
	"github.com/ManyakRus/starter/micro"
	"os"
)

// Getenv - возвращает переменную окружения
func Getenv(Name string, IsRequired bool) string {
	TextError := "Need fill OS environment variable: "
	Otvet := os.Getenv(Name)
	if IsRequired == true && Otvet == "" {
		log.Error(TextError + Name)
	}

	return Otvet
}

// Set_FieldFromEnv_String - устанавливает значение поля из переменной окружения
// Параметры:
// Object - указатель на структуру
// FieldName - имя поля
func Set_FieldFromEnv_String(StructReference any, FieldName string, IsRequired bool) {
	Value := Getenv(FieldName, IsRequired)

	err := micro.SetFieldValue(StructReference, FieldName, Value)

	if err != nil {
		err = fmt.Errorf("SetFieldFrom() FieldName: %s error: %w", FieldName, err)
		log.Error(err)
		return
	}
}
