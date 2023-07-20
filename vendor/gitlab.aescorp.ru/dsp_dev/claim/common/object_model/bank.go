package object_model

// Bank Банки (справочник).
type Bank struct {
	CommonStruct
	GroupStruct
	NameStruct
	OrganizationID       int64  `json:"organization_id"       gorm:"column:organization_id;default:null"`
	City                 string `json:"city_name"             gorm:"column:city_name;default:\"\""`
	BIK                  string `json:"bik"                   gorm:"column:bik;default:\"\""`
	CorrespondentAccount string `json:"correspondent_account" gorm:"column:correspondent_account;default:\"\""`
	ConnectionID         int64  `json:"connection_id"         gorm:"column:connection_id;default:null"`
}
