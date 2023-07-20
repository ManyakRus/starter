package object_model

import (
	"time"
)

type OrganizationCasebook struct {
	CommonStruct
	INN            string    `json:"inn"             gorm:"column:inn;default:\"\""`
	JSONFileID     int64     `json:"json_file_id"    gorm:"column:json_file_id;default:null"`
	KPP            string    `json:"kpp"             gorm:"column:kpp;default:\"\""`
	OrganizationID int64     `json:"organization_id" gorm:"column:organization_id;default:null"`
	PDFFileID      int64     `json:"pdf_file_id"     gorm:"column:pdf_file_id;default:null"`
	UpdatedAt      time.Time `json:"updated_at"      gorm:"column:updated_at;default:null"`
}
