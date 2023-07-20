package object_model

// GenderType Пол (справочник).
type GenderType struct {
	CommonStruct
	Name string `json:"name"            gorm:"column:name;default:\"\""`
}
