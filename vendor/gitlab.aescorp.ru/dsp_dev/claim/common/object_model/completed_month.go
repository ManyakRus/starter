package object_model

import (
	"time"
)

// CompletedMonth Закрытые месяцы
type CompletedMonth struct {
	CommonStruct
	ConnectionID     int64     `json:"connection_id" gorm:"column:connection_id;default:null"`
	AccountingAreaID int64     `json:"accounting_area_id" gorm:"column:accounting_area_id;default:null"`
	BillingMonth     time.Time `json:"billing_month" gorm:"column:billing_month;default:null"`
}
