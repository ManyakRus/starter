package object_model

import (
	"time"
)

// StateDuty Госпошлина
type StateDuty struct {
	CommonStruct
	GroupStruct
	NameStruct
	Sum           float64   `json:"sum"            gorm:"column:sum;default:null"`
	RequestNumber string    `json:"request_number" gorm:"column:request_number;default:\"\""`
	RequestDate   time.Time `json:"request_date"   gorm:"column:request_date;default:null"`
	CourtID       int64     `json:"court_id"       gorm:"column:court_id;default:null"`
	LawsuitID     int64     `json:"lawsuit_id"     gorm:"column:lawsuit_id;default:null"`
}
