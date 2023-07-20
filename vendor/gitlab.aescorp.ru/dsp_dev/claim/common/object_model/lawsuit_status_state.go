package object_model

// LawsuitStatusState История статусов дела.
type LawsuitStatusState struct {
	CommonStruct
	LawsuitID int64   `json:"lawsuit_id"      gorm:"column:lawsuit_id;default:null"`
	StatusID  int64   `json:"status_id"       gorm:"column:status_id;default:null"`
	Tag       string  `json:"tag"             gorm:"column:tag;default:\"\""`
	TotalDebt float64 `json:"total_debt"      gorm:"column:total_debt;default:null"`
}
