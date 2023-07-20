package object_model

// OrganizationStateType Состояние организации (справочник).
type OrganizationStateType struct {
	CommonStruct
	NameStruct
	Code               string `json:"code"                gorm:"column:code;default:\"1\""`
	ActionIndividual   string `json:"action_individual"   gorm:"action_individual:code;default:\"none\""`   // include exclude none
	ActionOrganization string `json:"action_organization" gorm:"action_organization:code;default:\"none\""` // include exclude none
	Color              string `json:"color"               gorm:"color:code;default:\"\""`                   // red yellow green
}
