package object_model

// DocumentLink Связи документов
type DocumentLink struct {
	CommonStruct
	ConnectionID  int64   `json:"connection_id" gorm:"column:connection_id;default:null"`
	ContractID    int64   `json:"contract_id"   gorm:"column:contract_id;default:null"`
	CorrectionSum float64 `json:"correction_sum" gorm:"column:correction_sum;not null;default:0"`
	Document1ID   int64   `json:"document1_id"  gorm:"column:document1_id;default:null"`
	Document2ID   int64   `json:"document2_id"  gorm:"column:document2_id;default:null"`
	LinkTypeID    int64   `json:"link_type_id"  gorm:"column:link_type_id;default:null"`
}
