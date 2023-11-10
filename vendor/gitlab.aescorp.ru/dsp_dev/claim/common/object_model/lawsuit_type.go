package object_model

// LawsuitType Типы дел (справочник).
type LawsuitType struct {
	CommonStruct
	NameStruct
	Code string `json:"code" gorm:"column:code;default:null"`
}
