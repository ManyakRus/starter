package object_model

import (
	"time"
)

// PaymentSchedule График платежей по договору
type PaymentSchedule struct {
	CommonStruct
	ConnectionID int64     `json:"connection_id"   gorm:"column:connection_id;default:null"`
	ContractID   int64     `json:"contract_id"     gorm:"column:contract_id;default:null"`
	DateFrom     time.Time `json:"date_from"       gorm:"column:date_from;default:null"`
	DateTo       time.Time `json:"date_to"         gorm:"column:date_to;default:null"`
	Day          int       `json:"day"             gorm:"column:day;not null"`
	Percent      int       `json:"percent"         gorm:"column:percent;not null"`
}
