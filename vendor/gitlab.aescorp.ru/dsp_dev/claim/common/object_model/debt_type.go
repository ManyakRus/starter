package object_model

// DebtType Виды задолженности
type DebtType struct {
	CommonStruct
	GroupStruct
	NameStruct
	CodeNSI      int   `json:"code_nsi"        gorm:"column:code_nsi;default:0"`
	ExtCode      int   `json:"ext_code"        gorm:"column:ext_code;default:0"`
	ConnectionID int64 `json:"connection_id"   gorm:"column:connection_id;default:null"`
}
