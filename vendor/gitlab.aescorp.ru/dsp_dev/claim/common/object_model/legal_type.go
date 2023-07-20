package object_model

// LegalType Тип юридического лица (справочник).
type LegalType struct {
	CommonStruct
	NameStruct
	IsIndividual bool `json:"is_individual" gorm:"column:is_individual;default:false"`
}
