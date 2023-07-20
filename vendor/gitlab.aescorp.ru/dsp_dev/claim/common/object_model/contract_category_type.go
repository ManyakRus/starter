package object_model

// ContractCategoryType Категория договоров (справочник).
type ContractCategoryType struct {
	CommonStruct
	NameStruct
	GroupStruct
	Code         string `json:"code"          gorm:"column:code;default:\"\""`
	ConnectionID int64  `json:"connection_id" gorm:"column:connection_id;default:null"`
}
