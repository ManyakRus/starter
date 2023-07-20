package object_model

type GroupStruct struct {
	IsGroup  bool  `json:"is_group"  gorm:"column:is_group;default:false"`
	ParentID int64 `json:"parent_id" gorm:"column:parent_id;default:null"`
}
