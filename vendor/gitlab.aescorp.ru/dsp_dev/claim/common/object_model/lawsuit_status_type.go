package object_model

// LawsuitStatusType Статусы дел (справочник).
type LawsuitStatusType struct {
	CommonStruct
	NameStruct
	Code string `json:"code" gorm:"column:code;default:0"`
}
