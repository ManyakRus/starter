// пакет для микрофункций с логгером

package microl

import (
	"fmt"
	"github.com/ManyakRus/starter/constants"
	"github.com/ManyakRus/starter/log"
	"github.com/ManyakRus/starter/micro"
	"os"
	"strconv"
	"time"
)

// Getenv - возвращает переменную окружения
func Getenv(Name string, IsRequired bool) string {
	TextError := "Need fill OS environment variable: "
	Otvet, IsFind := os.LookupEnv(Name)
	if IsFind == true {
		return Otvet
	}

	if IsRequired == true {
		log.Panic(TextError + Name)
	} else {
		log.Warn(TextError + Name)
	}

	return Otvet
}

// Set_FieldFromEnv_String - устанавливает значение поля из переменной окружения
// Параметры:
// Object - указатель на структуру
// FieldName - имя поля
// IsRequired - обязательное ли поле
func Set_FieldFromEnv_String(StructReference any, FieldName string, IsRequired bool) {
	Value := Getenv(FieldName, IsRequired)

	err := micro.SetFieldValue(StructReference, FieldName, Value)

	if err != nil {
		err = fmt.Errorf("SetFieldFrom() FieldName: %s error: %w", FieldName, err)
		log.Error(err)
		return
	}
}

// Set_FieldFromEnv_Int - устанавливает значение поля из переменной окружения
// Параметры:
// Object - указатель на структуру
// FieldName - имя поля
// IsRequired - обязательное ли поле
func Set_FieldFromEnv_Int(StructReference any, FieldName string, IsRequired bool) {
	sValue := Getenv(FieldName, IsRequired)

	Value, err := strconv.Atoi(sValue)
	if err != nil {
		err = fmt.Errorf("Atoi() FieldName: %s error: %w", FieldName, err)
		log.Error(err)
		return
	}

	err = micro.SetFieldValue(StructReference, FieldName, Value)

	if err != nil {
		err = fmt.Errorf("SetFieldFrom() FieldName: %s error: %w", FieldName, err)
		log.Error(err)
		return
	}
}

// Set_FieldFromEnv_Int64 - устанавливает значение поля из переменной окружения
// Параметры:
// Object - указатель на структуру
// FieldName - имя поля
// IsRequired - обязательное ли поле
func Set_FieldFromEnv_Int64(StructReference any, FieldName string, IsRequired bool) {
	sValue := Getenv(FieldName, IsRequired)

	Value, err := strconv.ParseInt(sValue, 10, 64)
	if err != nil {
		err = fmt.Errorf("ParseInt() FieldName: %s error: %w", FieldName, err)
		log.Error(err)
		return
	}

	err = micro.SetFieldValue(StructReference, FieldName, Value)

	if err != nil {
		err = fmt.Errorf("SetFieldFrom() FieldName: %s error: %w", FieldName, err)
		log.Error(err)
		return
	}
}

// Set_FieldFromEnv_Int32 - устанавливает значение поля из переменной окружения
// Параметры:
// Object - указатель на структуру
// FieldName - имя поля
// IsRequired - обязательное ли поле
func Set_FieldFromEnv_Int32(StructReference any, FieldName string, IsRequired bool) {
	sValue := Getenv(FieldName, IsRequired)

	var Value int32
	Value, err := micro.Int32FromString(sValue)
	if err != nil {
		err = fmt.Errorf("ParseInt() FieldName: %s error: %w", FieldName, err)
		log.Error(err)
		return
	}

	err = micro.SetFieldValue(StructReference, FieldName, Value)

	if err != nil {
		err = fmt.Errorf("SetFieldFrom() FieldName: %s error: %w", FieldName, err)
		log.Error(err)
		return
	}
}

// Set_FieldFromEnv_Time - устанавливает значение поля из переменной окружения
// Параметры:
// Object - указатель на структуру
// FieldName - имя поля
// IsRequired - обязательное ли поле
func Set_FieldFromEnv_Time(StructReference any, FieldName string, IsRequired bool) {
	sValue := Getenv(FieldName, IsRequired)

	Value, err := time.Parse(constants.LayoutDateTimeRus, sValue)
	if err != nil {
		err = fmt.Errorf("time.Parse() FieldName: %s error: %w", FieldName, err)
		log.Error(err)
		return
	}

	err = micro.SetFieldValue(StructReference, FieldName, Value)

	if err != nil {
		err = fmt.Errorf("SetFieldFrom() FieldName: %s error: %w", FieldName, err)
		log.Error(err)
		return
	}
}

// Set_FieldFromEnv_Date - устанавливает значение поля из переменной окружения
// Параметры:
// Object - указатель на структуру
// FieldName - имя поля
// IsRequired - обязательное ли поле
func Set_FieldFromEnv_Date(StructReference any, FieldName string, IsRequired bool) {
	sValue := Getenv(FieldName, IsRequired)

	Value, err := time.Parse(constants.LayoutDateRus, sValue)
	if err != nil {
		err = fmt.Errorf("time.Parse() FieldName: %s error: %w", FieldName, err)
		log.Error(err)
		return
	}

	err = micro.SetFieldValue(StructReference, FieldName, Value)

	if err != nil {
		err = fmt.Errorf("SetFieldFrom() FieldName: %s error: %w", FieldName, err)
		log.Error(err)
		return
	}
}

// Set_FieldFromEnv_Bool - устанавливает значение поля из переменной окружения
// Параметры:
// Object - указатель на структуру
// FieldName - имя поля
// IsRequired - обязательное ли поле
func Set_FieldFromEnv_Bool(StructReference any, FieldName string, IsRequired bool) {
	sValue := Getenv(FieldName, IsRequired)

	Value := micro.BoolFromString(sValue)

	err := micro.SetFieldValue(StructReference, FieldName, Value)

	if err != nil {
		err = fmt.Errorf("SetFieldFrom() FieldName: %s error: %w", FieldName, err)
		log.Error(err)
		return
	}
}

// ShowTimePassed - показывает время прошедшее с момента старта
// запускать:
// defer micro.ShowTimePassed(time.Now())
func ShowTimePassed(StartAt time.Time) {
	log.Debugf("Time passed: %s\n", time.Since(StartAt))
}

// ShowTimePassed_FormatText - показывает время прошедшее с момента старта
// запускать:
// defer micro.ShowTimePassed(time.Now())
func ShowTimePassed_FormatText(FormatText string, StartAt time.Time) {
	log.Debugf(FormatText, time.Since(StartAt))
}

// ShowTimePassedSeconds - показывает время секунд прошедшее с момента старта
// запускать:
// defer micro.ShowTimePassedSeconds(time.Now())
func ShowTimePassedSeconds(StartAt time.Time) {
	log.Debugf("Time passed: %s\n", time.Since(StartAt).Round(time.Second))
}

// ShowTimePassedMilliSeconds - показывает время миллисекунд прошедшее с момента старта
// запускать:
// defer micro.ShowTimePassedMilliSeconds(time.Now())
func ShowTimePassedMilliSeconds(StartAt time.Time) {
	log.Debugf("Time passed: %s\n", time.Since(StartAt).Round(time.Millisecond))
}

// Set_StructField - устанавливает значение поля из переменной окружения
// Параметры:
// Object - указатель на структуру
// FieldName - имя поля
// Value - значение
func Set_StructField(StructReference any, FieldName string, Value any) {
	err := micro.SetFieldValue(StructReference, FieldName, Value)

	if err != nil {
		err = fmt.Errorf("Set_StructField() FieldName: %s error: %w", FieldName, err)
		log.Error(err)
		return
	}
}
