package object_model

type ExtLinkStruct struct {
	TableNameID int64 `json:"table_name_id" gorm:"column:table_name_id;default:null"`
	TableRowID  int64 `json:"table_row_id"  gorm:"column:table_row_id;default:null"`
}
