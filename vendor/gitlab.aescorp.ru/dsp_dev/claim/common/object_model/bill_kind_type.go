package object_model

// BillKindType Вид платежа
type BillKindType struct {
	CommonStruct
	NameStruct
	Code int `json:"code" gorm:"column:code;default:null"`
}
