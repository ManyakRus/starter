package object_model

// DocumentLinkType Тип связи документов
type DocumentLinkType struct {
	CommonStruct
	NameStruct
	Code int `json:"code" gorm:"column:code;default:null"`
}
