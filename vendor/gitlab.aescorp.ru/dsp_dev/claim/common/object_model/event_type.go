package object_model

// EventType Типы событий (справочник).
type EventType struct {
	CommonStruct
	NameStruct
	Code string `json:"code" gorm:"column:code;default:null"`
}
