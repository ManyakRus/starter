package object_model

// DocumentType Тип документов (справочник).
type DocumentType struct {
	CommonStruct
	NameStruct
	IsService     bool   `json:"is_service"     gorm:"column:is_service;default:false"`
	IsVisible     bool   `json:"is_visible"     gorm:"column:is_visible;default:false"`
	ShortName     string `json:"short_name"     gorm:"column:short_name;default:\"\""`
	Type          int    `json:"type"           gorm:"column:type;default:0"`
	IncomeExpense int    `json:"income_expense" gorm:"column:income_expense;default:null"`
	ConnectionID  int64  `json:"connection_id"  gorm:"column:connection_id;default:null"`
}
