package object_model

// LawsuitReasonType Причина отбора для претензии (Справочник).
type LawsuitReasonType struct {
	CommonStruct
	NameStruct
	Code string `json:"code" gorm:"column:code;default:\"\""`
}
