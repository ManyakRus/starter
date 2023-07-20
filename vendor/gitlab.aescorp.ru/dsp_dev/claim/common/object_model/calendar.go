package object_model

import (
	"time"
)

// Calendar - список выходных
type Calendar struct {
	CommonStruct
	Date      time.Time `json:"date" gorm:"column:date;default:null"`
	Hours     int       `json:"hours" gorm:"column:hours;default:0"`
	Days      int       `json:"days" gorm:"column:days;default:0"`
	DayTypeID int64     `json:"day_type_id" gorm:"column:day_type_id;default:null"`
	Comment   string    `json:"comment" gorm:"column:comment;default:\"\""`
}
