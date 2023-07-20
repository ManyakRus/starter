package object_model

// AccountingArea Области учёта
type AccountingArea struct {
	CommonStruct
	NameStruct
	ConnectionID int64 `json:"connection_id" gorm:"column:connection_id;default:null"`
	Code         int   `json:"code"          gorm:"column:code;default:null"`
}
