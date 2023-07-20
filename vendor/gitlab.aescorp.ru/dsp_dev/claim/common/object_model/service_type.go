package object_model

// ServiceType Типы услуг (справочник).
type ServiceType struct {
	CommonStruct
	NameStruct
	GroupStruct
	Code              int    `json:"code"                gorm:"column:code;default:0"`
	FullName          string `json:"full_name"           gorm:"column:full_name;default:\"\""`
	Measure           string `json:"measure"             gorm:"column:measure;default:\"\""`
	ServiceProviderID int64  `json:"service_provider_id" gorm:"column:service_provider_id;default:null"`
	ConnectionID      int64  `json:"connection_id"       gorm:"column:connection_id;default:null"`
}
