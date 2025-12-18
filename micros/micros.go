// пакет для микрофункций с stopapp

package micros

import (
	"github.com/ManyakRus/starter/stopapp"
	"time"
)

// IEnable - интерфейс для включения
type IEnable interface {
	Enable()
}

// IDisable - интерфейс для отключения
type IDisable interface {
	Disable()
}

// EnableAfterDuration - выполняет Enable() после паузы
func EnableAfterDuration(Object IEnable, Duration time.Duration) {
	if Object == nil {
		return
	}
	stopapp.GetWaitGroup_Main().Add(1)
	go EnableAfterDuration_go(Object, Duration)
}

// EnableAfterMilliSeconds - выполняет Enable() после паузы
func EnableAfterMilliSeconds(Object IEnable, MilliSeconds int) {
	if Object == nil {
		return
	}
	stopapp.GetWaitGroup_Main().Add(1)
	go EnableAfterDuration_go(Object, time.Duration(MilliSeconds)*time.Millisecond)
}

// EnableAfterDuration_go - горутина, выполняет Enable() после паузы
func EnableAfterDuration_go(Object IEnable, Duration time.Duration) {
	defer stopapp.GetWaitGroup_Main().Done()

	if Object == nil {
		return
	}

	//
	time.Sleep(Duration)
	Object.Enable()

}
