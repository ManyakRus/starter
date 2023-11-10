package object_model

import (
	"time"
)

// LawsuitInvoice Счета фактуры относящиеся к делу.
type LawsuitInvoice struct {
	CommonStruct
	ClosedAt     time.Time `json:"closed_at"     gorm:"column:closed_at;default:null"`
	ClosedSum    float64   `json:"closed_sum"    gorm:"column:closed_sum;default:null"`
	Count        int64     `json:"count"         gorm:"column:count;not null"`
	Document     Document  `json:"document"      gorm:"-:all"`
	DocumentID   int64     `json:"document_id"   gorm:"column:document_id;default:null"`        // Document
	DocumentSum  float64   `json:"document_sum"  gorm:"column:document_sum;not null;default:0"` // Сумма указанная в платёжном документе
	IsClosed     bool      `json:"is_closed"     gorm:"is_closed:tag;default:false"`
	IsCorrective bool      `json:"is_corrective" gorm:"column:is_corrective;default:false"`
	LawsuitID    int64     `json:"lawsuit_id"    gorm:"column:lawsuit_id;default:null"` // Lawsuit
	Sum          float64   `json:"sum"           gorm:"column:sum;not null;default:0"`  // Сумма фактуры после коррекции
}
