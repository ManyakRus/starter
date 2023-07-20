package object_model

import (
	"time"
)

// ===========================================================================
// ===== Объекты =====
// ===========================================================================

type CommonStruct struct {
	ID         int64     `json:"id"          gorm:"column:id;primaryKey;autoIncrement:true"`
	ExtID      int64     `json:"ext_id"      gorm:"column:ext_id;default:null"`
	CreatedAt  time.Time `json:"created_at"  gorm:"column:created_at;autoCreateTime"`
	ModifiedAt time.Time `json:"modified_at" gorm:"column:modified_at;autoUpdateTime"`
	DeletedAt  time.Time `json:"deleted_at"  gorm:"column:deleted_at;default:null"`
	IsDeleted  bool      `json:"is_deleted"  gorm:"column:is_deleted;default:false"`
}
