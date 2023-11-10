package object_model

import "github.com/google/uuid"

// ChannelType Тип канала (справочник).
type ChannelType struct {
	CommonStruct
	NameStruct
	NotifierID uuid.UUID `json:"notifier_id"            gorm:"type:uuid;column:notifier_id;default:\"\""` //ИД как в Notifier
	Code       int       `json:"code"        gorm:"column:code;default:0"`
}
