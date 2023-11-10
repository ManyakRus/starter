package object_model

import (
	"gitlab.aescorp.ru/dsp_dev/claim/common/object_model/alias"
	"gitlab.aescorp.ru/dsp_dev/claim/common/object_model/front/front_format"
)

// LawsuitPayment -- платежи относящиеся к делу
type LawsuitPayment struct {
	CommonStruct
	ID           alias.PaymentId `json:"id"          gorm:"column:id;primaryKey;autoIncrement:true"`
	Document     Document        `json:"document"     gorm:"-:all"`
	DocumentID   int64           `json:"document_id"  gorm:"column:document_id;default:null"`        // Document
	DocumentSum  float64         `json:"document_sum" gorm:"column:document_sum;not null;default:0"` // Сумма указанная в платёжном документе
	Invoice      LawsuitInvoice  `json:"invoice"      gorm:"-:all"`
	InvoiceID    alias.InvoiceId `json:"invoice_id"   gorm:"column:invoice_id;default:null"` // LawsuitInvoice
	IsCorrective bool            `json:"is_corrective"    gorm:"column:is_corrective;default:false"`
	LawsuitID    alias.LawsuitId `json:"lawsuit_id"   gorm:"column:lawsuit_id;default:null"` // Lawsuit
	Sum          float64         `json:"sum"          gorm:"column:sum;not null;default:0"`  // Сумма погашения после коррекции
}

// // IsAfterNotify -- возвращает признак создания платежа после уведомления
// func (sf *LawsuitPayment) IsAfterNotify(contractNumber alias.ContractNumber) alias.IsAfterNotify {
// 	lawsuit := NewLawsuit(contractNumber)
// 	controlDate := lawsuit.CreatedAt
// 	strControlDate := controlDate.Local().Format("2006-01-02 15:04:05")
// 	paymentDate := sf.CreatedAt
// 	strPaymentDate := paymentDate.Local().Format("2006-01-02 15:04:05")
// 	return strControlDate < strPaymentDate
// }

// RegisteredAt -- омент регистрации платежа в системе
func (sf *LawsuitPayment) RegisteredAt() alias.PaymentRegisteredAt {
	strDate := front_format.FrontTime(sf.CreatedAt)
	return alias.PaymentRegisteredAt(strDate)
}

// DatePayAt -- возвращает момент оплаты
func (sf *LawsuitPayment) DatePayAt() alias.FrontDate {
	frontDate := front_format.FrontDate(sf.Document.DocumentAt)
	return frontDate
}

// InvoiceId -- возвращает ID привязанной С/Ф
func (sf *LawsuitPayment) InvoiceId() alias.InvoiceId {
	return sf.InvoiceID
}

// Id -- возвращает ID платёжки
func (sf *LawsuitPayment) Id() alias.PaymentId {
	return sf.ID
}
