// пакет для микрофункций с логгером

package microl

import (
	"github.com/ManyakRus/starter/log"
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
