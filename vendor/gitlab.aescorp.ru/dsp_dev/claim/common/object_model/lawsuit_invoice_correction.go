package object_model

// LawsuitInvoiceCorrection -
type LawsuitInvoiceCorrection struct {
	CommonStruct
	CorrectionDocumentID  int64   `json:"correction_document_id" gorm:"column:correction_document_id;default:null"`
	CorrectionDocumentSum float64 `json:"correction_document_sum" gorm:"column:correction_document_sum;default:null"`
	InvoiceDocumentID     int64   `json:"invoice_document_id"    gorm:"column:invoice_document_id;default:null"`
	LawsuitID             int64   `json:"lawsuit_id"             gorm:"column:lawsuit_id;default:null"` // Lawsuit
}
