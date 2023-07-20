package object_model

// ClaimType Типы исков (справочник).
type ClaimType struct {
	CommonStruct
	NameStruct
	Code string `json:"code" gorm:"column:code;default:\"\""`
}
