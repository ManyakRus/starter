package object_model

// ===========================================================================
// ===== Списки =====
// ===========================================================================

type NameStruct struct {
	Description string `json:"description" gorm:"column:description;default:\"\""`
	Name        string `json:"name"        gorm:"column:name;default:\"\""`
}
