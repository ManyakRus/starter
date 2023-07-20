package object_model

import (
	"time"
)

// Document Документ.
type Document struct {
	CommonStruct // Id с/ф в СТЕК; Дата формирования
	GroupStruct
	Analytics      string    `json:"analytics"        gorm:"column:analytics;default:\"\""`    // Тип начисления (окончательный, пени, ограничения, по суду);
	Balance        float64   `json:"balance"          gorm:"column:balance;default:null"`      // Неоплаченный остаток С/ф
	BillKindID     int64     `json:"bill_kind_id"     gorm:"column:bill_kind_id;default:null"` //
	BillingMonth   time.Time `json:"billing_month"    gorm:"column:billing_month;default:null"`
	ConnectionID   int64     `json:"connection_id"    gorm:"column:connection_id;default:null"`
	ContractID     int64     `json:"contract_id"      gorm:"column:contract_id;default:null"` // Номер договора;
	Count          int64     `json:"count"            gorm:"column:count;not null"`           // Квт.ч в С/ф;
	DebtSum        float64   `json:"debt_sum"         gorm:"column:debt_sum;default:null"`    // Долг в претензии
	DocumentAt     time.Time `json:"document_at"      gorm:"column:document_at;default:null"` // Дата С/ф;
	DocumentSum    float64   `json:"document_sum"     gorm:"column:document_sum;not null"`    // Начислено по С/ф
	DocumentTypeID int64     `json:"document_type_id" gorm:"column:document_type_id;default:null"`
	Note           string    `json:"note"             gorm:"column:note;default:\"\""`         // Примечание, в частности назначение платежа
	Number         string    `json:"number"           gorm:"column:number;default:\"\""`       // Номер С/ф;
	NumberFull     string    `json:"number_full"      gorm:"column:number_full;default:\"\""`  // Полный номер С/ф;
	PayDeadline    time.Time `json:"pay_deadline"     gorm:"column:pay_deadline;default:null"` // День когда уже пошла просрочка
	PayFrom        time.Time `json:"pay_from"         gorm:"column:pay_from;default:null"`     // Период С/ф; с
	PayTo          time.Time `json:"pay_to"           gorm:"column:pay_to;default:null"`       // Период С/ф; по
	Payment        float64   `json:"payment"          gorm:"column:payment;default:null"`      // Оплата по С/ф
	Reason         string    `json:"reason"           gorm:"column:reason;default:\"\""`
	ReversalID     int64     `json:"reversal_id"      gorm:"column:reversal_id;default:null"` // Указатель на исправленный документ
}
