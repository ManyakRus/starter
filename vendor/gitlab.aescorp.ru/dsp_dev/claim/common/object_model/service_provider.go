package object_model

import (
	"time"
)

// ServiceProvider Поставщик услуг (справочник).
type ServiceProvider struct {
	CommonStruct
	OrganizationID int64     `json:"organization_id" gorm:"column:organization_id;default:null"`
	ConnectionID   int64     `json:"connection_id"   gorm:"column:connection_id;default:null"`
	DateFrom       time.Time `json:"date_from"       gorm:"column:date_from;default:null"`
	DateTo         time.Time `json:"date_to"         gorm:"column:date_to;default:null"`
}
