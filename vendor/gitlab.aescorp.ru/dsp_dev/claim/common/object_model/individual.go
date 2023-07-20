package object_model

import (
	"time"
)

// Individual Физическое лицо (справочник).
type Individual struct {
	CommonStruct
	NameStruct
	BirthDate    time.Time `json:"birth_date"      gorm:"column:birth_date;default:null"`
	DeathDate    time.Time `json:"death_date"      gorm:"column:death_date;default:null"`
	Email        string    `json:"email"           gorm:"column:email;default:\"\""`
	GenderID     int64     `json:"gender_id"       gorm:"column:gender_id;default:null"`
	INN          string    `json:"inn"             gorm:"column:inn;default:\"\""`
	ParentName   string    `json:"parent_name"     gorm:"column:parent_name;default:\"\""`
	Phone        string    `json:"phone"           gorm:"column:phone;default:\"\""`
	SNILS        string    `json:"snils"           gorm:"column:snils;default:\"\""`
	SecondName   string    `json:"second_name"     gorm:"column:second_name;default:\"\""`
	ConnectionID int64     `json:"connection_id"   gorm:"column:connection_id;default:null"`
}
