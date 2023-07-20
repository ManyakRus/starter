package object_model

// LawsuitStageType Этапы дел (справочник).
type LawsuitStageType struct {
	CommonStruct
	NameStruct
	Code string `json:"code" gorm:"column:code;default:0"`
}
