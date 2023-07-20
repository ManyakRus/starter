package object_model

import (
	"time"
)

// Balance Сальдо договора.
type Balance struct {
	CommonStruct
	BillingMonth      time.Time `json:"billing_month"       gorm:"column:billing_month;default:null"`
	ConnectionID      int64     `json:"connection_id"       gorm:"column:connection_id;default:null"`
	ContractID        int64     `json:"contract_id"         gorm:"column:contract_id;default:null"`
	DebtTypeID        int64     `json:"debt_type_id"        gorm:"column:debt_type_id;default:null"`
	DocumentAt        time.Time `json:"document_at"         gorm:"column:document_at;default:null"`
	DocumentInvoiceID int64     `json:"document_invoice_id" gorm:"column:document_invoice_id;default:null"`
	DocumentPaymentID int64     `json:"document_payment_id" gorm:"column:document_payment_id;default:null"`
	Sum               float64   `json:"sum"                 gorm:"column:sum;default:null"`
}
