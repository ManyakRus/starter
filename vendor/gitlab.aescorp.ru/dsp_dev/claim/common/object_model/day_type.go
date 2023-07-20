package object_model

// DayType - тип рабочего дня
type DayType struct {
	CommonStruct
	NameStruct
	ShortName string `json:"short_name" gorm:"column:short_name;default:\"\""`
	IsWorkDay bool   `json:"is_work_day" gorm:"column:is_work_day;default:false"`
}
