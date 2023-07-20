package object_model

// Court Суды (справочник).
type Court struct {
	CommonStruct
	GroupStruct
	NameStruct
	OrganizationID int64  `json:"organization_id" gorm:"column:organization_id;default:null"`
	City           string `json:"city_name"       gorm:"column:city_name;default:\"\""`
}
