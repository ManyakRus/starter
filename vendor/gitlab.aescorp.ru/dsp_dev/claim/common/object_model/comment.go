package object_model

// Comment Комментарии.
type Comment struct {
	CommonStruct
	ExtLinkStruct
	Message string `json:"message"       gorm:"column:message;default:\"\""`
}
