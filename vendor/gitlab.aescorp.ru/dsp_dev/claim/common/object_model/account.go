package object_model

// Account Лицевой счёт
type Account struct {
	CommonStruct
	Bank   Bank   `json:"bank"    gorm:"-:all"`
	BankID int64  `json:"bank_id" gorm:"column:bank_id;default:null"`
	Number string `json:"number"  gorm:"column:number;default:\"\""`
}
