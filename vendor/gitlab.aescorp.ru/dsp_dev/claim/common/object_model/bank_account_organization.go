package object_model

// BankAccountOrganization Соответствие банка - лицевого счёта - организации.
type BankAccountOrganization struct {
	CommonStruct
	Bank           Bank         `json:"bank"            gorm:"-:all"`
	BankID         int64        `json:"bank_id"         gorm:"column:bank_id;default:null"`
	ConnectionID   int64        `json:"connection_id"   gorm:"column:connection_id;default:null"`
	AccountNumber  string       `json:"account_number"  gorm:"column:account_number;default:\"\""`
	Organization   Organization `json:"organization"    gorm:"-:all"`
	OrganizationID int64        `json:"organization_id" gorm:"column:organization_id;default:null"`
}
