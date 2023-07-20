package object_model

// Message Сообщения (входящие и исходящие).
type Message struct {
	CommonStruct

	// Тип сообщения
	// Идентификатор сообщения
	// Дата отправки
	// Статус отправки
	// Код доставки
	// Статус доставки

	ChannelTypeID   int64  `json:"channel_type_id"   gorm:"column:channel_type_id;default:null"`
	Code            string `json:"code"              gorm:"column:code;default:\"\""`
	Data            string `json:"data"              gorm:"column:data;default:\"\""`
	DirectionTypeID int64  `json:"direction_type_id" gorm:"column:direction_type_id;default:null"`
	LawsuitID       int64  `json:"lawsuit_id"        gorm:"column:lawsuit_id;default:null"`
	Result          string `json:"result"            gorm:"column:result;default:\"\""`
}
