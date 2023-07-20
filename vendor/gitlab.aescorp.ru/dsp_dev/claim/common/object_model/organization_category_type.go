package object_model

// OrganizationCategoryType Категория организаций (справочник).
type OrganizationCategoryType struct {
	CommonStruct
	NameStruct
	GroupStruct
	ConnectionID int64 `json:"connection_id" gorm:"column:connection_id;default:null"`
}
