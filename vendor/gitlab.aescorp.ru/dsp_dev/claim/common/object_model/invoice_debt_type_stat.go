package object_model

import (
	"gitlab.aescorp.ru/dsp_dev/claim/common/object_model/alias"
	"time"
)

type InvoiceDebtTypeStats struct {
	CommonStruct
	DebtTypeID int64           `json:"debt_type_id"  gorm:"column:debt_type_id;default:null"`
	InvoiceID  alias.InvoiceId `json:"invoice_id"    gorm:"column:invoice_id;default:null"`
	StateAt    time.Time       `json:"state_at"      gorm:"column:state_at;default:null"`
}
