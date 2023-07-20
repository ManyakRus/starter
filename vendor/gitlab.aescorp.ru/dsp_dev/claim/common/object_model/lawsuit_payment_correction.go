package object_model

// LawsuitPaymentCorrection -
type LawsuitPaymentCorrection struct {
	CommonStruct
	CorrectionDocumentID  int64   `json:"correction_document_id" gorm:"column:correction_document_id;default:null"`
	CorrectionDocumentSum float64 `json:"correction_document_sum" gorm:"column:correction_document_sum;default:null"`
	LawsuitID             int64   `json:"lawsuit_id"             gorm:"column:lawsuit_id;default:null"` // Lawsuit
	PaymentDocumentID     int64   `json:"payment_document_id"    gorm:"column:payment_document_id;default:null"`
}
