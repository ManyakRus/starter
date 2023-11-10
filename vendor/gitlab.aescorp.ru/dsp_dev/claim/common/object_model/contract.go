package object_model

import (
	"encoding/json"
	"reflect"
	"time"

	"github.com/vmihailenco/msgpack/v5"
	"gitlab.aescorp.ru/dsp_dev/claim/common/object_model/alias"
)

// VersionStructContract - версия структуры модели, с учётом имен и типов полей
var versionStructContract uint32

// Contract Договоры.
type Contract struct {
	CommonStruct
	GroupStruct
	BeginAt                 time.Time            `json:"begin_at"               gorm:"column:begin_at"`
	BranchID                int64                `json:"branch_id"              gorm:"column:branch_id;default:null"`
	Category                ContractCategoryType `json:"category"               gorm:"-:all"`
	CategoryID              int64                `json:"category_id"            gorm:"column:category_id;default:null"`
	ConnectionID            int64                `json:"connection_id"          gorm:"column:connection_id;default:null"`
	CuratorClaim            Employee             `json:"curator_claim"          gorm:"-:all"`
	CuratorClaimID          int64                `json:"curator_claim_id"       gorm:"column:curator_claim_id;default:null"`
	CuratorContract         Employee             `json:"curator_contract"       gorm:"-:all"`
	CuratorContractID       int64                `json:"curator_contract_id"    gorm:"column:curator_contract_id;default:null"`
	CuratorLegal            Employee             `json:"curator_legal"          gorm:"-:all"`
	CuratorLegalID          int64                `json:"curator_legal_id"       gorm:"column:curator_legal_id;default:null"`
	CuratorPayment          Employee             `json:"curator_payment"        gorm:"-:all"`
	CuratorPaymentID        int64                `json:"curator_payment_id"     gorm:"column:curator_payment_id;default:null"`
	CuratorTechAudit        Employee             `json:"curator_tech_audit"     gorm:"-:all"`
	CuratorTechAuditID      int64                `json:"curator_tech_audit_id"  gorm:"column:curator_tech_audit_id;default:null"`
	DaysToResolveClaim      int                  `json:"days_to_resolve_claim"  gorm:"column:days_to_resolve_claim"`
	Description             string               `json:"description"            gorm:"column:description;default:\"\""`
	Email                   string               `json:"email"                  gorm:"column:email;default:\"\""`
	EndAt                   time.Time            `json:"end_at"                 gorm:"column:end_at"`
	ErrorFromStackAt        time.Time            `json:"error_from_stack_at"    gorm:"column:error_from_stack_at"`
	IndividualID            int64                `json:"individual_id"          gorm:"column:individual_id;default:null"`
	IsErrorFromStack        bool                 `json:"is_error_from_stack"    gorm:"column:is_error_from_stack;default:false"`
	IsIndOrganization       bool                 `json:"is_ind_organization"    gorm:"column:is_ind_organization;default:false"`
	IsOrganization          bool                 `json:"is_organization"        gorm:"column:is_organization;default:false"`
	IsValidEmail            bool                 `json:"is_valid_email"         gorm:"column:is_valid_email;default:true"`
	Number                  alias.ContractNumber `json:"number"                 gorm:"column:number;default:\"\""`
	Organization            Organization         `json:"organization"           gorm:"-:all"`
	OrganizationConsigneeID int64                `json:"organization_consignee_id" gorm:"column:organization_consignee_id;default:null"` // Грузополучатель - consignee
	OrganizationCustomerID  int64                `json:"organization_customer_id"  gorm:"column:organization_customer_id;default:null"`  // Заказчик - customer
	OrganizationPayerID     int64                `json:"organization_payer_id"     gorm:"column:organization_payer_id;default:null"`     // Плательщик - payer
	PostAddress             string               `json:"post_address"           gorm:"column:post_address;default:\"\""`
	SignAt                  time.Time            `json:"sign_at"                gorm:"column:sign_at"`
	Status                  string               `json:"status"                 gorm:"column:status;default:\"\""`
	TerminateAt             time.Time            `json:"terminate_at"           gorm:"column:terminate_at"`
	PaymentDays             []PaymentDay         `json:"payment_days"`      // Дни платежей
	PaymentSchedules        []PaymentSchedule    `json:"payment_schedules"` // Графики платежей
}

// NewContract Договор
func NewContract() Contract {
	return Contract{}
}

func AsContract(b []byte) (Contract, error) {
	c := NewContract()
	err := msgpack.Unmarshal(b, &c)
	if err != nil {
		return NewContract(), err
	}
	return c, nil
}

func ContractAsBytes(c *Contract) ([]byte, error) {
	b, err := msgpack.Marshal(c)
	if err != nil {
		return nil, err
	}
	return b, nil
}

// GetStructVersion - возвращает версию модели
func (c Contract) GetStructVersion() uint32 {
	if versionStructContract == 0 {
		versionStructContract = CalcStructVersion(reflect.TypeOf(c))
	}

	return versionStructContract
}

// GetModelFromJSON - заполняет Модель из строки JSON
func (c *Contract) GetModelFromJSON(sModel string) error {
	var err error

	var bytes []byte
	bytes = []byte(sModel)

	err = json.Unmarshal(bytes, c)

	return err
}

// GetJSON - возвращает строку json из модели
func (c Contract) GetJSON() (string, error) {
	var ReturnVar string
	var err error

	bytes, err := json.Marshal(c)
	if err != nil {
		return ReturnVar, err
	}
	ReturnVar = string(bytes)
	return ReturnVar, err
}
