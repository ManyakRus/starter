package object_model

// Branch Филиалы (справочник).
type Branch struct {
	CommonStruct
	GroupStruct
	NameStruct
	Code             int64  `json:"code"            gorm:"column:code;default:null"`
	OrganizationID   int64  `json:"organization_id" gorm:"column:organization_id;default:null"`
	PersonalAreaLink string `json:"personal_area_link" gorm:"personal_area_link:tag;default:\"\""`
}
