package object_model

import (
	"fmt"
	"strings"
	"time"

	"github.com/vmihailenco/msgpack/v5"
	"gitlab.aescorp.ru/dsp_dev/claim/common/object_model/object_view"
)

// regexp  ищет все типы string с default null  string.{1,}.{1,}gorm.{1,}default:null
// regexp  ищет все типы bool с default null  bool.{1,}.{1,}gorm.{1,}default:null
// regexp  ищет все типы time.Time с default null  time.Time.{1,}gorm.{1,}default:null

// ===========================================================================
// ===== Списки =====
// ===========================================================================

// TypeRef общие, как правило, фиксированные справочники
type TypeRef struct {
	ChannelTypes              []ChannelType              `json:"channel_types"`               // 1. Каналы сообщений
	ClaimTypes                []ClaimType                `json:"claim_types"`                 // 2. Типы дел
	ContractCategoryTypes     []ContractCategoryType     `json:"contract_category_types"`     // 3. * Категории договоров
	DebtTypes                 []DebtType                 `json:"debt_types"`                  // 4. * Типы задолженностей
	DirectionTypes            []DirectionType            `json:"direction_types"`             // 5. Направления сообщений
	DocumentLinkTypes         []DocumentLinkType         `json:"document_link_types"`         // 6. * Типы связей документов
	DocumentTypes             []DocumentType             `json:"document_types"`              // 7. * Типы документов
	FileTypes                 []FileType                 `json:"file_types"`                  // 8. Типы файлов
	GenderTypes               []GenderType               `json:"gender_types"`                // 9. * Пол
	LawsuitReasonTypes        []LawsuitReasonType        `json:"lawsuit_reason_types"`        // 10. * Причина отбора для претензии (Справочник).
	LawsuitStageTypes         []LawsuitStageType         `json:"lawsuit_stage_types"`         // 11. * Этапы дел (справочник).
	LawsuitStatusTypes        []LawsuitStatusType        `json:"lawsuit_status_types"`        // 12. * Статусы дел (справочник).
	LegalTypes                []LegalType                `json:"legal_types"`                 // 13. Тип юридического лица
	OrganizationCategoryTypes []OrganizationCategoryType `json:"organization_category_types"` // 14. Категории организаций
	ServiceTypes              []ServiceType              `json:"service_types"`               // 15. Типы услуг
	TableNames                []TableName                `json:"table_names"`                 // 16. * Имена таблиц для привязок
	WhiteListReasonTypes      []WhiteListReasonType      `json:"white_list_reason_types"`     // 17. * Причина добавления в белый список
}

// CommonRef Справочники
type CommonRef struct {
	Banks            []Bank              `json:"banks"`             // 1. Банки
	Branches         []Branch            `json:"branches"`          // 2. Отделения
	Courts           []Court             `json:"courts"`            // 3. Суды
	FileTemplates    []FileTemplate      `json:"file_templates"`    // 4. Шаблоны документов
	ServiceProviders []ServiceProvider   `json:"service_providers"` // 5. Поставщики услуг
	UserRoles        []UserRole          `json:"user_roles"`        // 6. Роли сотрудников
	WhiteList        []ContractWhiteItem `json:"white_list"`        // 7. Белый список
}

// BriefCase Набор данных для конкретного портфеля
type BriefCase struct {
	Lawsuit      Lawsuit              // Дело
	ChangeItems  []ChangeItem         `json:"change_items"`  // 3. История изменений
	Comments     []Comment            `json:"comments"`      // 4. Комментарии
	Files        []File               `json:"files"`         // 7. Файлы
	Invoices     []LawsuitInvoice     `json:"invoices"`      // 8. Счета фактуры
	Messages     []Message            `json:"messages"`      // 9. Сообщения
	Payments     []LawsuitPayment     `json:"payments"`      // 10. Платежи
	StateDuties  []StateDuty          `json:"state_duties"`  // 11. Гос.пошлина
	StatusStates []LawsuitStatusState `json:"status_states"` // 12. История статусов дела
	// TODO Добавить период претензии
}

// ClaimWork ПИР
type ClaimWork struct {
	BriefCases []BriefCase `json:"brief_cases"`
}

// ===========================================================================
// ===== Объекты =====
// ===========================================================================

type CommonStruct struct {
	ID         int64     `json:"id"          gorm:"column:id;primaryKey;autoIncrement:true"`
	ExtID      int64     `json:"ext_id"      gorm:"column:ext_id;default:null"`
	CreatedAt  time.Time `json:"created_at"  gorm:"column:created_at;autoCreateTime"`
	ModifiedAt time.Time `json:"modified_at" gorm:"column:modified_at;autoUpdateTime"`
	DeletedAt  time.Time `json:"deleted_at"  gorm:"column:deleted_at;default:null"`
	IsDeleted  bool      `json:"is_deleted"  gorm:"column:is_deleted;default:false"`
}

type NameStruct struct {
	Description string `json:"description" gorm:"column:description;default:\"\""`
	Name        string `json:"name"        gorm:"column:name;default:\"\""`
}

type GroupStruct struct {
	IsGroup  bool  `json:"is_group"  gorm:"column:is_group;default:false"`
	ParentID int64 `json:"parent_id" gorm:"column:parent_id;default:null"`
}

type ExtLinkStruct struct {
	TableNameID int64 `json:"table_name_id" gorm:"column:table_name_id;default:null"`
	TableRowID  int64 `json:"table_row_id"  gorm:"column:table_row_id;default:null"`
}

// ===========================================================================

// Contract Договоры.
type Contract struct {
	CommonStruct
	GroupStruct
	BeginAt            time.Time         `json:"begin_at"               gorm:"column:begin_at"`
	BranchID           int64             `json:"branch_id"              gorm:"column:branch_id;default:null"`
	CategoryID         int64             `json:"category_id"            gorm:"column:category_id;default:null"`
	ConnectionID       int64             `json:"connection_id"          gorm:"column:connection_id;default:null"`
	CuratorClaim       Employee          `json:"curator_claim"          gorm:"-:all"`
	CuratorClaimID     int64             `json:"curator_claim_id"       gorm:"column:curator_claim_id;default:null"`
	CuratorContract    Employee          `json:"curator_contract"       gorm:"-:all"`
	CuratorContractID  int64             `json:"curator_contract_id"    gorm:"column:curator_contract_id;default:null"`
	CuratorLegal       Employee          `json:"curator_legal"          gorm:"-:all"`
	CuratorLegalID     int64             `json:"curator_legal_id"       gorm:"column:curator_legal_id;default:null"`
	CuratorPayment     Employee          `json:"curator_payment"        gorm:"-:all"`
	CuratorPaymentID   int64             `json:"curator_payment_id"     gorm:"column:curator_payment_id;default:null"`
	CuratorTechAudit   Employee          `json:"curator_tech_audit"     gorm:"-:all"`
	CuratorTechAuditID int64             `json:"curator_tech_audit_id"  gorm:"column:curator_tech_audit_id;default:null"`
	DaysToResolveClaim int               `json:"days_to_resolve_claim"  gorm:"column:days_to_resolve_claim"`
	Description        string            `json:"description"            gorm:"column:description;default:\"\""`
	Email              string            `json:"email"                  gorm:"column:email;default:\"\""`
	EndAt              time.Time         `json:"end_at"                 gorm:"column:end_at"`
	IndividualID       int64             `json:"individual_id"          gorm:"column:individual_id;default:null"`
	IsIndOrganization  bool              `json:"is_ind_organization"    gorm:"column:is_ind_organization;default:false"`
	IsOrganization     bool              `json:"is_organization"        gorm:"column:is_organization;default:false"`
	IsValidEmail       bool              `json:"is_valid_email"         gorm:"column:is_valid_email;default:false"`
	Number             string            `json:"number"                 gorm:"column:number;default:\"\""`
	Organization       Organization      `json:"organization"           gorm:"-:all"`
	OrganizationID     int64             `json:"organization_id"        gorm:"column:organization_id;default:null"`
	PostAddress        string            `json:"post_address"           gorm:"column:post_address;default:\"\""`
	SignAt             time.Time         `json:"sign_at"                gorm:"column:sign_at"`
	Status             string            `json:"status"                 gorm:"column:status;default:\"\""`
	TerminateAt        time.Time         `json:"terminate_at"           gorm:"column:terminate_at"`
	PaymentDays        []PaymentDay      `json:"payment_days"`      // Дни платежей
	PaymentSchedules   []PaymentSchedule `json:"payment_schedules"` // Графики платежей
}

// LawsuitStatusType Статусы дел (справочник).
type LawsuitStatusType struct {
	CommonStruct
	NameStruct
	Code string `json:"code" gorm:"column:code;default:\"\""`
}

// LawsuitStageType Этапы дел (справочник).
type LawsuitStageType struct {
	CommonStruct
	NameStruct
	Code string `json:"code" gorm:"column:code;default:\"\""`
}

// LawsuitReasonType Причина отбора для претензии (Справочник).
type LawsuitReasonType struct {
	CommonStruct
	NameStruct
	Code string `json:"code" gorm:"column:code;default:\"\""`
}

// Lawsuit Дело. Объединяет весь набор данных по конкретному должнику.
type Lawsuit struct {
	CommonStruct
	GroupStruct
	NameStruct
	Branch                    Branch            `json:"branch"                  gorm:"-:all"`
	BranchID                  int64             `json:"branch_id"               gorm:"column:branch_id;default:null"`
	Chance                    string            `json:"chance"                  gorm:"column:chance;default:\"\""`
	ClaimAt                   time.Time         `json:"claim_at"                gorm:"column:claim_at;default:null"` // Уведомление о задолженности. Оплатить до.
	ClaimPeriodStr            string            `json:"claim_period_str"        gorm:"column:claim_period_str;default:\"\""`
	ClaimType                 ClaimType         `json:"claim_type"              gorm:"-:all"` // Тип задолженности
	ClaimTypeID               int64             `json:"claim_type_id"           gorm:"column:claim_type_id;default:null"`
	ClosedAt                  time.Time         `json:"closed_at"               gorm:"column:closed_at;default:null"`
	Contract                  Contract          `json:"contract"                gorm:"-:all"` // Договор
	ContractID                int64             `json:"contract_id"             gorm:"column:contract_id;default:null"`
	ControlledAt              time.Time         `json:"controlled_at"           gorm:"column:controlled_at;default:null"`
	Court                     Court             `json:"court"                   gorm:"-:all"`
	CourtID                   int64             `json:"court_id"                gorm:"column:court_id;default:null"`
	CuratorClaim              Employee          `json:"curator_claim"           gorm:"-:all"`
	CuratorClaimID            int64             `json:"curator_claim_id"        gorm:"column:curator_claim_id;default:null"`
	CuratorContract           Employee          `json:"curator_contract"        gorm:"-:all"`
	CuratorContractID         int64             `json:"curator_contract_id"     gorm:"column:curator_contract_id;default:null"`
	CuratorLegal              Employee          `json:"curator_legal"           gorm:"-:all"`
	CuratorLegalID            int64             `json:"curator_legal_id"        gorm:"column:curator_legal_id;default:null"`
	CuratorPayment            Employee          `json:"curator_payment"         gorm:"-:all"`
	CuratorPaymentID          int64             `json:"curator_payment_id"      gorm:"column:curator_payment_id;default:null"`
	CuratorTechAudit          Employee          `json:"curator_tech_audit"      gorm:"-:all"`
	CuratorTechAuditID        int64             `json:"curator_tech_audit_id"   gorm:"column:curator_tech_audit_id;default:null"`
	DateFrom                  time.Time         `json:"date_from"               gorm:"column:date_from;default:null"`
	DateTo                    time.Time         `json:"date_to"                 gorm:"column:date_to;default:null"`
	DebtSum                   float64           `json:"debt_sum"                gorm:"column:debt_sum;default:0"`                             // Общая сумма долга
	InvoiceSum                float64           `json:"invoice_sum"             gorm:"column:invoice_sum;default:0"`                          // Сумма долга за период
	NotifyClaimAt             time.Time         `json:"notify_claim_at"         gorm:"column:notify_claim_at;default:null"`                   // Уведомление о задолженности. Дата отправки.
	NotifyClaimChannel        int               `json:"notify_claim_channel"    gorm:"column:notify_claim_channel;default:null"`              // Уведомление о задолженности. Канал отправки.
	NotifyClaimCode           int               `json:"notify_claim_code"       gorm:"column:notify_claim_code;default:null"`                 // Уведомление о задолженности. Код доставки из НСИ.
	NotifyClaimDone           bool              `json:"notify_claim_done"       gorm:"column:notify_claim_done;default:null"`                 // Уведомление о задолженности. Факт отправки.
	NotifyClaimMailingCode    string            `json:"notify_claim_mailing_code" gorm:"column:notify_claim_mailing_code;default:null"`       // Уведомление о задолженности. Уникальный код отправки.
	NotifyPretrialAt          time.Time         `json:"notify_pretrial_at"      gorm:"column:notify_pretrial_at;default:null"`                // Досудебная претензия. Дата отправки.
	NotifyPretrialChannel     int               `json:"notify_pretrial_channel" gorm:"column:notify_pretrial_channel;default:null"`           // Досудебная претензия. Канал отправки.
	NotifyPretrialCode        int               `json:"notify_pretrial_code"    gorm:"column:notify_pretrial_code;default:null"`              // Досудебная претензия. Код доставки из НСИ.
	NotifyPretrialDone        bool              `json:"notify_pretrial_done"    gorm:"column:notify_pretrial_done;default:null"`              // Досудебная претензия. Факт отправки.
	NotifyPretrialMailingCode string            `json:"notify_pretrial_mailing_code" gorm:"column:notify_pretrial_mailing_code;default:null"` // Досудебная претензия. Уникальный код отправки.
	Number                    string            `json:"number"                  gorm:"column:number;default:\"\""`
	NumberClaim               string            `json:"number_claim"            gorm:"column:number_claim;default:\"\""`
	NumberTrial               string            `json:"number_trial"            gorm:"column:number_trial;default:\"\""`
	PaySum                    float64           `json:"pay_sum"                 gorm:"column:pay_sum;default:0"` // Платежи
	Penalty                   float64           `json:"penalty"                 gorm:"column:penalty;default:0"`
	Penny                     float64           `json:"penny"                   gorm:"column:penny;default:0"`
	Percent317                float64           `json:"percent_317"             gorm:"column:percent_317;default:0"`
	Percent395                float64           `json:"percent_395"             gorm:"column:percent_395;default:0"`
	PretrialAt                time.Time         `json:"pretrial_at"             gorm:"column:pretrial_at;default:null"` // Досудебная претензия. Оплатить до.
	ProcessStartedAt          time.Time         `json:"process_started_at"      gorm:"column:process_started_at;default:null"`
	Reason                    LawsuitReasonType `json:"reason"                  gorm:"-:all"`
	ReasonID                  int64             `json:"reason_id"               gorm:"column:reason_id;default:null"`
	Stage                     LawsuitStageType  `json:"stage"                   gorm:"-:all"` // Этап
	StageAt                   time.Time         `json:"stage_at"                gorm:"column:stage_at;default:null"`
	StageID                   int64             `json:"stage_id"                gorm:"column:stage_id;default:null"`
	StateDuty                 float64           `json:"state_duty"              gorm:"column:state_duty;default:0"` // Пошлина
	Status                    LawsuitStatusType `json:"status"                  gorm:"-:all"`                       // Статус
	StatusAt                  time.Time         `json:"status_at"               gorm:"column:status_at;default:null"`
	StatusID                  int64             `json:"status_id"               gorm:"column:status_id;default:null"`
	Tag                       string            `json:"tag"                     gorm:"column:tag;default:\"\""`
	UnknownPayments           bool              `json:"unknown_payments"        gorm:"unknown_payments:tag;default:false"` // "c не разнесёнными платежами"
}

// LawsuitInvoice Счета фактуры относящиеся к делу.
type LawsuitInvoice struct {
	CommonStruct
	LawsuitID  int64     `json:"lawsuit_id"  gorm:"column:lawsuit_id;default:null"`  // Lawsuit
	DocumentID int64     `json:"document_id" gorm:"column:document_id;default:null"` // Document
	Document   Document  `json:"document"    gorm:"-:all"`
	Sum        float64   `json:"sum"         gorm:"column:sum;default:null"`
	Count      int64     `json:"count"       gorm:"column:count;not null"`
	IsClosed   bool      `json:"is_closed"   gorm:"is_closed:tag;default:false"`
	ClosedAt   time.Time `json:"closed_at"   gorm:"column:closed_at;default:null"`
}

// LawsuitPayment Платежи относящиеся к делу.
type LawsuitPayment struct {
	CommonStruct
	LawsuitID  int64          `json:"lawsuit_id"  gorm:"column:lawsuit_id;default:null"`  // Lawsuit
	DocumentID int64          `json:"document_id" gorm:"column:document_id;default:null"` // Document
	Document   Document       `json:"document"    gorm:"-:all"`
	InvoiceID  int64          `json:"invoice_id"  gorm:"column:invoice_id;default:null"` // LawsuitInvoice
	Invoice    LawsuitInvoice `json:"invoice"     gorm:"-:all"`
	Sum        float64        `json:"sum"         gorm:"column:sum;default:null"`
}

// LawsuitStatusState История статусов дела.
type LawsuitStatusState struct {
	CommonStruct
	LawsuitID int64   `json:"lawsuit_id"      gorm:"column:lawsuit_id;default:null"`
	StatusID  int64   `json:"status_id"       gorm:"column:status_id;default:null"`
	Tag       string  `json:"tag"             gorm:"column:tag;default:\"\""`
	TotalDebt float64 `json:"total_debt"      gorm:"column:total_debt;default:null"`
}

// TODO LawsuitPeriod Период претензии

// FileType Тип файла (справочник).
type FileType struct {
	CommonStruct
	NameStruct
}

// StateDuty Госпошлина
type StateDuty struct {
	CommonStruct
	GroupStruct
	NameStruct
	Sum           float64   `json:"sum"            gorm:"column:sum;default:null"`
	RequestNumber string    `json:"request_number" gorm:"column:request_number;default:\"\""`
	RequestDate   time.Time `json:"request_date"   gorm:"column:request_date;default:null"`
	CourtID       int64     `json:"court_id"       gorm:"column:court_id;default:null"`
	LawsuitID     int64     `json:"lawsuit_id"     gorm:"column:lawsuit_id;default:null"`
}

// DebtType Виды задолженности
type DebtType struct {
	CommonStruct
	GroupStruct
	NameStruct
	ExtCode      int   `json:"ext_code"        gorm:"column:ext_code;default:0"`
	ConnectionID int64 `json:"connection_id"   gorm:"column:connection_id;default:null"`
}

// FileTemplate Шаблоны файлов (справочник).
type FileTemplate struct {
	CommonStruct
	NameStruct
	FileID string `json:"file_id" gorm:"column:file_id;default:\"\""`
}

// Message Сообщения (входящие и исходящие).
type Message struct {
	CommonStruct
	ChannelTypeID   int64  `json:"channel_type_id"   gorm:"column:channel_type_id;default:null"`
	Code            string `json:"code"              gorm:"column:code;default:\"\""`
	Data            string `json:"data"              gorm:"column:data;default:\"\""`
	DirectionTypeID int64  `json:"direction_type_id" gorm:"column:direction_type_id;default:null"`
	LawsuitID       int64  `json:"lawsuit_id"        gorm:"column:lawsuit_id;default:null"`
	Result          string `json:"result"            gorm:"column:result;default:\"\""`
}

// Employee Сотрудники (Справочник).
type Employee struct {
	CommonStruct
	NameStruct
	GroupStruct
	BranchID     int64  `json:"branch_id"       gorm:"column:branch_id;default:null"`
	Email        string `json:"email"           gorm:"column:email;default:\"\""`
	IsActive     bool   `json:"is_active"       gorm:"column:is_active;default:false"`
	Login        string `json:"login"           gorm:"column:login;default:\"\""`
	ParentName   string `json:"parent_name"     gorm:"column:parent_name;default:\"\""`
	Phone        string `json:"phone"           gorm:"column:phone;default:\"\""`
	Photo        string `json:"photo"           gorm:"column:photo;default:\"\""`
	Position     string `json:"position"        gorm:"column:position;default:\"\""`
	SecondName   string `json:"second_name"     gorm:"column:second_name;default:\"\""`
	Tag          string `json:"tag"             gorm:"column:tag;default:\"\""`
	ConnectionID int64  `json:"connection_id"   gorm:"column:connection_id;default:null"`
}

// Branch Филиалы (справочник).
type Branch struct {
	CommonStruct
	GroupStruct
	NameStruct
	Code             int64  `json:"code"            gorm:"column:code;default:null"`
	OrganizationID   int64  `json:"organization_id" gorm:"column:organization_id;default:null"`
	PersonalAreaLink string `json:"personal_area_link" gorm:"personal_area_link:tag;default:\"\""`
}

// Organization Юридическое лицо (справочник).
type Organization struct {
	CommonStruct
	NameStruct
	GroupStruct
	IsBankrupt     bool      `json:"is_bankrupt"     gorm:"column:is_bankrupt;default:false"`
	BankruptAt     time.Time `json:"bankrupt_at"     gorm:"column:bankrupt_at"`
	BookkeeperName string    `json:"bookkeeper_name" gorm:"column:bookkeeper_name;default:\"\""`
	CategoryID     int64     `json:"category_id"     gorm:"column:category_id;default:null"`
	ConnectionID   int64     `json:"connection_id"   gorm:"column:connection_id;default:null"`
	Email          string    `json:"email"           gorm:"column:email;default:\"\""`
	FullName       string    `json:"full_name"       gorm:"column:full_name;default:\"\""`
	INN            string    `json:"inn"             gorm:"column:inn;default:\"\""`
	IsActive       bool      `json:"is_active"       gorm:"column:is_active;default:false"`
	KPP            string    `json:"kpp"             gorm:"column:kpp;default:\"\""`
	LegalAddress   string    `json:"legal_address"   gorm:"column:legal_address;default:\"\""`
	LiquidateAt    time.Time `json:"liquidate_at"    gorm:"column:liquidate_at"`
	IsLiquidated   bool      `json:"is_liquidated"   gorm:"column:is_liquidated;default:false"`
	ManagerName    string    `json:"manager_name"    gorm:"column:manager_name;default:\"\""`
	OGRN           string    `json:"ogrn"            gorm:"column:ogrn;default:\"\""`
	OKATO          string    `json:"okato"           gorm:"column:okato;default:\"\""`
	OKPO           string    `json:"okpo"            gorm:"column:okpo;default:\"\""`
	Phone          string    `json:"phone"           gorm:"column:phone;default:\"\""`
	PostAddress    string    `json:"post_address"    gorm:"column:post_address;default:\"\""`
	RegistrationAt time.Time `json:"registration_at" gorm:"column:registration_at;default:null"`
	WWW            string    `json:"www"             gorm:"column:www;default:\"\""`
	Accounts       []Account `json:"accounts"        gorm:"-:all"`
}

// OrganizationCategoryType Категория организаций (справочник).
type OrganizationCategoryType struct {
	CommonStruct
	NameStruct
	GroupStruct
	ConnectionID int64 `json:"connection_id" gorm:"column:connection_id;default:null"`
}

// ContractCategoryType Категория договоров (справочник).
type ContractCategoryType struct {
	CommonStruct
	NameStruct
	GroupStruct
	Code         string `json:"code"          gorm:"column:code;default:\"\""`
	ConnectionID int64  `json:"connection_id" gorm:"column:connection_id;default:null"`
}

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

// GenderType Пол (справочник).
type GenderType struct {
	CommonStruct
	Name string `json:"name"            gorm:"column:name;default:\"\""`
}

// DirectionType Направление передачи сообщения (справочник).
type DirectionType struct {
	CommonStruct
	NameStruct
}

// LegalType Тип юридического лица (справочник).
type LegalType struct {
	CommonStruct
	NameStruct
}

// ChannelType Тип канала (справочник).
type ChannelType struct {
	CommonStruct
	NameStruct
}

// ClaimType Типы исков (справочник).
type ClaimType struct {
	CommonStruct
	NameStruct
	Code string `json:"code" gorm:"column:code;default:\"\""`
}

// Document Документ.
type Document struct {
	CommonStruct // Id с/ф в СТЕК; Дата формирования
	GroupStruct
	Analytics      string    `json:"analytics"        gorm:"column:analytics;default:\"\""` // Тип начисления (окончательный, пени, ограничения, по суду);
	Balance        float64   `json:"balance"          gorm:"column:balance;default:null"`   // Неоплаченный остаток С/ф
	BillingMonth   time.Time `json:"billing_month"    gorm:"column:billing_month;default:null"`
	ConnectionID   int64     `json:"connection_id"    gorm:"column:connection_id;default:null"`
	ContractID     int64     `json:"contract_id"      gorm:"column:contract_id;default:null"` // Номер договора;
	Count          int64     `json:"count"            gorm:"column:count;not null"`           // Квт.ч в С/ф;
	DebtSum        float64   `json:"debt_sum"         gorm:"column:debt_sum;default:null"`    // Долг в претензии
	DocumentAt     time.Time `json:"document_at"      gorm:"column:document_at;default:null"` // Дата С/ф;
	DocumentSum    float64   `json:"document_sum"     gorm:"column:document_sum;not null"`    // Начислено по С/ф
	DocumentTypeID int64     `json:"document_type_id" gorm:"column:document_type_id;default:null"`
	Number         string    `json:"number"           gorm:"column:number;default:\"\""`   // Номер С/ф;
	PayFrom        time.Time `json:"pay_from"         gorm:"column:pay_from;default:null"` // Период С/ф; с
	PayTo          time.Time `json:"pay_to"           gorm:"column:pay_to;default:null"`   // Период С/ф; по
	Payment        float64   `json:"payment"          gorm:"column:payment;default:null"`  // Оплата по С/ф
	Reason         string    `json:"reason"           gorm:"column:reason;default:\"\""`
	ReversalID     int64     `json:"reversal_id"      gorm:"column:reversal_id;default:null"` // Указатель на исправленный документ
	Note           string    `json:"note"             gorm:"column:note;default:\"\""`        // Примечание, в частности назначение платежа
}

// Balance Сальдо договора.
type Balance struct {
	CommonStruct
	BillingMonth      time.Time `json:"billing_month"       gorm:"column:billing_month;default:null"`
	ConnectionID      int64     `json:"connection_id"       gorm:"column:connection_id;default:null"`
	ContractID        int64     `json:"contract_id"         gorm:"column:contract_id;default:null"`
	DebtTypeID        int64     `json:"debt_type_id"        gorm:"column:debt_type_id;default:null"`
	DocumentAt        time.Time `json:"document_at"         gorm:"column:document_at;default:null"`
	DocumentInvoiceID int64     `json:"document_invoice_id" gorm:"column:document_invoice_id;default:null"`
	DocumentPaymentID int64     `json:"document_payment_id" gorm:"column:document_payment_id;default:null"`
	Sum               float64   `json:"sum"                 gorm:"column:sum;default:null"`
}

// DocumentLink Связи документов
type DocumentLink struct {
	CommonStruct
	ConnectionID int64 `json:"connection_id" gorm:"column:connection_id;default:null"`
	ContractID   int64 `json:"contract_id"   gorm:"column:contract_id;default:null"`
	Document1ID  int64 `json:"document1_id"  gorm:"column:document1_id;default:null"`
	Document2ID  int64 `json:"document2_id"  gorm:"column:document2_id;default:null"`
	LinkTypeID   int64 `json:"link_type_id"  gorm:"column:link_type_id;default:null"`
}

// DocumentLinkType Тип связи документов
type DocumentLinkType struct {
	CommonStruct
	NameStruct
	Code int `json:"code" gorm:"column:code;default:null"`
}

// DocumentType Тип документов (справочник).
type DocumentType struct {
	CommonStruct
	NameStruct
	IsService     bool   `json:"is_service"     gorm:"column:is_service;default:false"`
	IsVisible     bool   `json:"is_visible"     gorm:"column:is_visible;default:false"`
	ShortName     string `json:"short_name"     gorm:"column:short_name;default:\"\""`
	Type          int    `json:"type"           gorm:"column:type;default:0"`
	IncomeExpense int    `json:"income_expense" gorm:"column:income_expense;default:null"`
	ConnectionID  int64  `json:"connection_id"  gorm:"column:connection_id;default:null"`
}

// UserRole Роль (справочник).
type UserRole struct {
	CommonStruct
	NameStruct
}

// PaymentSchedule График платежей по договору
type PaymentSchedule struct {
	CommonStruct
	ConnectionID int64     `json:"connection_id"   gorm:"column:connection_id;default:null"`
	ContractID   int64     `json:"contract_id"     gorm:"column:contract_id;default:null"`
	DateFrom     time.Time `json:"date_from"       gorm:"column:date_from;default:null"`
	DateTo       time.Time `json:"date_to"         gorm:"column:date_to;default:null"`
	Day          int       `json:"day"             gorm:"column:day;not null"`
	Percent      int       `json:"percent"         gorm:"column:percent;not null"`
}

// PaymentDay День платежа по договору
type PaymentDay struct {
	CommonStruct
	ConnectionID int64     `json:"connection_id"   gorm:"column:connection_id;default:null"`
	ContractID   int64     `json:"contract_id"     gorm:"column:contract_id;default:null"`
	DateFrom     time.Time `json:"date_from"       gorm:"column:date_from;default:null"`
	DateTo       time.Time `json:"date_to"         gorm:"column:date_to;default:null"`
	Day          int       `json:"day"             gorm:"column:day;not null"`
}

// ServiceProvider Поставщик услуг (справочник).
type ServiceProvider struct {
	CommonStruct
	OrganizationID int64     `json:"organization_id" gorm:"column:organization_id;default:null"`
	ConnectionID   int64     `json:"connection_id"   gorm:"column:connection_id;default:null"`
	DateFrom       time.Time `json:"date_from"       gorm:"column:date_from;default:null"`
	DateTo         time.Time `json:"date_to"         gorm:"column:date_to;default:null"`
}

// ServiceType Типы услуг (справочник).
type ServiceType struct {
	CommonStruct
	NameStruct
	GroupStruct
	Code              int    `json:"code"                gorm:"column:code;default:0"`
	FullName          string `json:"full_name"           gorm:"column:full_name;default:\"\""`
	Measure           string `json:"measure"             gorm:"column:measure;default:\"\""`
	ServiceProviderID int64  `json:"service_provider_id" gorm:"column:service_provider_id;default:null"`
	ConnectionID      int64  `json:"connection_id"       gorm:"column:connection_id;default:null"`
}

// Bank Банки (справочник).
type Bank struct {
	CommonStruct
	GroupStruct
	NameStruct
	OrganizationID       int64  `json:"organization_id"       gorm:"column:organization_id;default:null"`
	City                 string `json:"city_name"             gorm:"column:city_name;default:\"\""`
	BIK                  string `json:"bik"                   gorm:"column:bik;default:\"\""`
	CorrespondentAccount string `json:"correspondent_account" gorm:"column:correspondent_account;default:\"\""`
	ConnectionID         int64  `json:"connection_id"         gorm:"column:connection_id;default:null"`
}

// Court Суды (справочник).
type Court struct {
	CommonStruct
	GroupStruct
	NameStruct
	OrganizationID int64  `json:"organization_id" gorm:"column:organization_id;default:null"`
	City           string `json:"city_name"       gorm:"column:city_name;default:\"\""`
}

// File Файлы.
type File struct {
	CommonStruct
	GroupStruct
	NameStruct
	ExtLinkStruct
	BranchID   int64  `json:"branch_id"       gorm:"column:branch_id;default:null"`
	EmployeeID int64  `json:"employee_id"     gorm:"column:employee_id;default:null"`
	Extension  string `json:"extension"       gorm:"column:extension;default:\"\""`
	FileID     string `json:"file_id"         gorm:"column:file_id;default:\"\""`
	FileName   string `json:"file_name"       gorm:"column:file_name;default:\"\""`
	FileTypeID int64  `json:"file_type_id"    gorm:"column:file_type_id;default:null"`
	FullName   string `json:"full_name"       gorm:"column:full_name;default:\"\""`
	Size       int64  `json:"size"            gorm:"column:size;default:null"`
	TemplateID int64  `json:"template_id"     gorm:"column:template_id;default:null"`
	Version    int    `json:"version"         gorm:"column:version;default:0"`
}

// Comment Комментарии.
type Comment struct {
	CommonStruct
	ExtLinkStruct
	Message string `json:"message"       gorm:"column:message;default:\"\""`
}

// ChangeItem Изменения
// Key - изменённое поле
// Value - новое значение
// Prev - прежнее значение
type ChangeItem struct {
	CommonStruct
	ExtLinkStruct
	// TODO UserID
	// TODO Action
	// TODO Table
	// TODO Field
	Key   string `json:"key"   gorm:"column:key;default:\"\""`
	Value string `json:"value" gorm:"column:value;default:\"\""`
	Prev  string `json:"prev"  gorm:"column:prev;default:\"\""`
}

// TableName объект позволяющий привязать такие компоненты как комментарий или файл к разным объектам
type TableName struct {
	CommonStruct
	NameStruct
}

type Connection struct {
	ID       int64  `json:"id"        gorm:"column:id;primaryKey;autoIncrement:true"`
	Name     string `json:"name"      gorm:"column:name;default:\"\""`
	IsLegal  bool   `json:"is_legal"  gorm:"column:is_legal;default:false"`
	BranchId int64  `json:"branch_id" gorm:"column:branch_id;default:0"`
	Server   string `json:"server"    gorm:"column:server;default:\"\""`
	Port     string `json:"port"      gorm:"column:port;default:\"\""`
	DbName   string `json:"db_name"   gorm:"column:db_name;default:\"\""`
	DbScheme string `json:"db_scheme" gorm:"column:db_scheme;default:\"\""`
	Login    string `json:"login"     gorm:"column:login;default:\"\""`
	Password string `json:"password"  gorm:"column:password;default:\"\""`
}

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

// Account Лицевой счёт
type Account struct {
	CommonStruct
	Bank   Bank   `json:"bank"    gorm:"-:all"`
	BankID int64  `json:"bank_id" gorm:"column:bank_id;default:null"`
	Number string `json:"number"  gorm:"column:number;default:\"\""`
}

// Facsimile Соответствие участков ответственных и договоров
type Facsimile struct {
	CommonStruct
	Branch      string `json:"branch"      gorm:"column:branch;default:\"\""`
	Department  string `json:"department"  gorm:"column:department;default:\"\""`
	Responsible string `json:"responsible" gorm:"column:responsible;default:\"\""`
	Contract    string `json:"contract"    gorm:"column:contract;default:\"\""`
}

// AccountingArea Области учёта
type AccountingArea struct {
	CommonStruct
	NameStruct
	ConnectionID int64 `json:"connection_id" gorm:"column:connection_id;default:null"`
	Code         int   `json:"code"          gorm:"column:code;default:null"`
}

// CompletedMonth Закрытые месяцы
type CompletedMonth struct {
	CommonStruct
	ConnectionID     int64     `json:"connection_id" gorm:"column:connection_id;default:null"`
	AccountingAreaID int64     `json:"accounting_area_id" gorm:"column:accounting_area_id;default:null"`
	BillingMonth     time.Time `json:"billing_month" gorm:"column:billing_month;default:null"`
}

// ContractWhiteItem "Белый" список договоров. Кому не предъявляется претензия.
type ContractWhiteItem struct {
	CommonStruct
	Contract       Contract            `json:"contract"        gorm:"-:all"`
	ContractID     int64               `json:"contract_id"     gorm:"column:contract_id;default:null"`
	ContractNumber string              `json:"contract_number" gorm:"column:contract_number;default:null"`
	CreatedBy      Employee            `json:"created_by"      gorm:"-:all"`
	CreatedByID    int64               `json:"created_by_id"   gorm:"column:created_by_id;default:null"`
	DateFrom       time.Time           `json:"date_from"       gorm:"column:date_from;default:null"`
	DateTo         time.Time           `json:"date_to"         gorm:"column:date_to;default:null"`
	EDMSLink       string              `json:"edms_link"       gorm:"column:edms_link;default:\"\""`
	ModifiedBy     Employee            `json:"modified_by"     gorm:"-:all"`
	ModifiedByID   int64               `json:"modified_by_id"  gorm:"column:modified_by_id;default:null"`
	Note           string              `json:"note"            gorm:"column:note;default:\"\""`
	Reason         WhiteListReasonType `json:"reason"          gorm:"-:all"`
	ReasonID       int64               `json:"reason_id"       gorm:"column:reason_id;default:null"`
}

// WhiteListReasonType Причина добавления договора в "белый" список (справочник).
type WhiteListReasonType struct {
	CommonStruct
	NameStruct
	Code int `json:"code" gorm:"column:code;default:null"`
}

// ===========================================================================
// ===== Отчёты =====
// ===========================================================================

// ReportSummary Суммарный отчёт (дашборд)
type ReportSummary struct {
	// Всего претензий сформировано с начала года
	ClaimsTotal    int    `json:"claims_total"      gorm:"-:all"`
	ClaimsTotalSum string `json:"claims_total_sum"  gorm:"-:all"`

	// В работе (исключая завершённые статусы)
	ClaimsInWork    int    `json:"claims_in_work"      gorm:"-:all"`
	ClaimsInWorkSum string `json:"claims_in_work_sum"  gorm:"-:all"`

	// Завершено на этапе претензии в связи с оплатой с начала года
	ClaimsStatus2    int    `json:"claims_status_2"      gorm:"-:all"`
	ClaimsStatus2Sum string `json:"claims_status_2_sum"  gorm:"-:all"`

	// Ожидает передачи в исковое производство
	ClaimsStatus7    int    `json:"claims_status_7"      gorm:"-:all"`
	ClaimsStatus7Sum string `json:"claims_status_7_sum"  gorm:"-:all"`

	// Досудебная претензия на стадии мониторинга с направлением по e-mail
	ClaimsChannel1401    int    `json:"claims_channel_1401"      gorm:"-:all"`
	ClaimsChannel1401Sum string `json:"claims_channel_1401_sum"  gorm:"-:all"`

	// Досудебная претензия на стадии мониторинга с направлением заказным письмом
	ClaimsChannel1406    int    `json:"claims_channel_1406"      gorm:"-:all"`
	ClaimsChannel1406Sum string `json:"claims_channel_1406_sum"  gorm:"-:all"`

	// Претензии с платежами не разнесенными на счет-фактуры
	ClaimsWithUnknown    int    `json:"claims_with_unknown"      gorm:"-:all"`
	ClaimsWithUnknownSum string `json:"claims_with_unknown_sum"  gorm:"-:all"`
}

// ===========================================================================
// ===== Методы =====
// ===========================================================================

// NewTypeRef Новый набор справочников типов
func NewTypeRef() TypeRef {
	return TypeRef{}
}

func AsTypeRef(b []byte) (TypeRef, error) {
	c := NewTypeRef()
	err := msgpack.Unmarshal(b, &c)
	if err != nil {
		return NewTypeRef(), err
	}
	return c, nil
}

func TypeRefAsBytes(c *TypeRef) ([]byte, error) {
	b, err := msgpack.Marshal(c)
	if err != nil {
		return nil, err
	}
	return b, nil
}

// ===========================================================================

// NewCommonRef Новый набор глобальных справочников
func NewCommonRef() CommonRef {
	return CommonRef{}
}

func AsCommonRef(b []byte) (CommonRef, error) {
	c := NewCommonRef()
	err := msgpack.Unmarshal(b, &c)
	if err != nil {
		return NewCommonRef(), err
	}
	return c, nil
}

func CommonRefAsBytes(c *CommonRef) ([]byte, error) {
	b, err := msgpack.Marshal(c)
	if err != nil {
		return nil, err
	}
	return b, nil
}

// ===========================================================================

// NewLawsuit Новый объект дела
func NewLawsuit(number string) Lawsuit {
	return Lawsuit{
		Number:      fmt.Sprintf("ПФ_%s_%s", time.Now().Format("200601-02"), number),
		NumberClaim: fmt.Sprintf("ПР_%s_%s", time.Now().Format("200601-02"), number),
		NumberTrial: fmt.Sprintf("ПИ_%s_%s", time.Now().Format("200601-02"), number),
	}
}

func AsLawsuit(b []byte) (Lawsuit, error) {
	c := Lawsuit{}
	err := msgpack.Unmarshal(b, &c)
	if err != nil {
		return Lawsuit{}, err
	}
	return c, nil
}

func LawsuitAsBytes(c *Lawsuit) ([]byte, error) {
	b, err := msgpack.Marshal(c)
	if err != nil {
		return nil, err
	}
	return b, nil
}

// NewBriefCase Новый объект портфеля
func NewBriefCase() BriefCase {
	return BriefCase{}
}

func AsBriefCase(b []byte) (BriefCase, error) {
	c := NewBriefCase()
	err := msgpack.Unmarshal(b, &c)
	if err != nil {
		return NewBriefCase(), err
	}
	return c, nil
}

func BriefCaseAsBytes(c *BriefCase) ([]byte, error) {
	b, err := msgpack.Marshal(c)
	if err != nil {
		return nil, err
	}
	return b, nil
}

// ===========================================================================

// NewClaimWork Новый объект верхнего уровня, по сути содержит список портфелей.
func NewClaimWork() ClaimWork {
	return ClaimWork{}
}

func AsClaimWork(b []byte) (ClaimWork, error) {
	c := NewClaimWork()
	err := msgpack.Unmarshal(b, &c)
	if err != nil {
		return NewClaimWork(), err
	}
	return c, nil
}

func ClaimWorkAsBytes(c *ClaimWork) ([]byte, error) {
	b, err := msgpack.Marshal(c)
	if err != nil {
		return nil, err
	}
	return b, nil
}

// ===========================================================================

// NewFile Файл, который физически хранится в файловом хранилище
func NewFile() File {
	return File{}
}

func AsFile(b []byte) (File, error) {
	f := NewFile()
	err := msgpack.Unmarshal(b, &f)
	if err != nil {
		return NewFile(), err
	}
	return f, nil
}

func FileAsBytes(f *File) ([]byte, error) {
	b, err := msgpack.Marshal(f)
	if err != nil {
		return nil, err
	}
	return b, nil
}

// ===========================================================================

// NewEmployee Сотрудник
func NewEmployee() Employee {
	return Employee{}
}

func AsEmployee(b []byte) (Employee, error) {
	e := NewEmployee()
	err := msgpack.Unmarshal(b, &e)
	if err != nil {
		return NewEmployee(), err
	}
	return e, nil
}

func EmployeeAsBytes(e *Employee) ([]byte, error) {
	b, err := msgpack.Marshal(e)
	if err != nil {
		return nil, err
	}
	return b, nil
}

// ===========================================================================

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

// ===========================================================================

// NewReportSummary Суммарный отчёт (дашборд)
func NewReportSummary() ReportSummary {
	return ReportSummary{}
}

func AsReportSummary(b []byte) (ReportSummary, error) {
	c := NewReportSummary()
	err := msgpack.Unmarshal(b, &c)
	if err != nil {
		return NewReportSummary(), err
	}
	return c, nil
}

func ReportSummaryAsBytes(c *ReportSummary) ([]byte, error) {
	b, err := msgpack.Marshal(c)
	if err != nil {
		return nil, err
	}
	return b, nil
}

// ===========================================================================

// NewFacsimile Данные факсимиле
func NewFacsimile() Facsimile {
	return Facsimile{}
}

func AsFacsimile(b []byte) (Facsimile, error) {
	c := NewFacsimile()
	err := msgpack.Unmarshal(b, &c)
	if err != nil {
		return NewFacsimile(), err
	}
	return c, nil
}

func FacsimileAsBytes(c *Facsimile) ([]byte, error) {
	b, err := msgpack.Marshal(c)
	if err != nil {
		return nil, err
	}
	return b, nil
}

// ===========================================================================

// BriefCaseView выборка
func BriefCaseView(bc *BriefCase, c *CommonRef, t *TypeRef, useFormat bool) (map[string]interface{}, error) {
	result := make(map[string]interface{})

	if bc == nil {
		return result, fmt.Errorf("MailTemplateView, BriefCase is nil")
	}
	if c == nil {
		return result, fmt.Errorf("MailTemplateView, CommonRef is nil")
	}
	if t == nil {
		return result, fmt.Errorf("MailTemplateView, TypeRef is nil")
	}

	lawSuit := bc.Lawsuit

	// ID
	result["Lawsuit_ID"] = lawSuit.ID
	// Дата претензии
	if useFormat {
		result["Lawsuit_CreatedAt"] = formatDate(lawSuit.CreatedAt)
	} else {
		result["Lawsuit_CreatedAt"] = lawSuit.CreatedAt
	}
	// Номер претензии
	result["Lawsuit_Number"] = lawSuit.Number
	result["Lawsuit_NumberClaim"] = lawSuit.NumberClaim
	result["Lawsuit_NumberTrial"] = lawSuit.NumberTrial

	// TODO View LawsuitStageTypes Улучшить поиск
	stage := "Неизвестно"
	for i := 0; i < len(t.LawsuitStageTypes); i++ {
		if t.LawsuitStageTypes[i].ID == lawSuit.StageID {
			stage = t.LawsuitStageTypes[i].Name
			break
		}
	}
	// Этап - для фильтрации
	result["Lawsuit_StageID"] = lawSuit.StageID
	// Этап - для вывода в таблицу
	result["Lawsuit_Stage"] = stage
	// Дата установки этапа
	if useFormat {
		result["Lawsuit_StageAt"] = formatDate(lawSuit.StageAt)
	} else {
		result["Lawsuit_StageAt"] = lawSuit.StageAt
	}

	// TODO View LawsuitStatusTypes Улучшить поиск
	status := "Неизвестно"
	for i := 0; i < len(t.LawsuitStatusTypes); i++ {
		if t.LawsuitStatusTypes[i].ID == lawSuit.StatusID {
			status = t.LawsuitStatusTypes[i].Name
			break
		}
	}
	// Статус - для фильтрации
	result["Lawsuit_StatusID"] = lawSuit.StatusID
	// Статус - для вывода в таблицу
	result["Lawsuit_Status"] = status
	// Дата установки статуса
	if useFormat {
		result["Lawsuit_StatusAt"] = formatDate(lawSuit.StatusAt)
	} else {
		result["Lawsuit_StatusAt"] = lawSuit.StatusAt
	}

	// TODO View LawsuitReasonTypes Улучшить поиск
	reason := "Неизвестно"
	for i := 0; i < len(t.LawsuitReasonTypes); i++ {
		if t.LawsuitReasonTypes[i].ID == lawSuit.ReasonID {
			reason = t.LawsuitReasonTypes[i].Name
			break
		}
	}
	// Статус - для фильтрации
	result["Lawsuit_ReasonID"] = lawSuit.ReasonID
	// Статус - для вывода в таблицу
	result["Lawsuit_Reason"] = reason

	// TODO View ClaimTypes Улучшить поиск
	claimType := "Неизвестно"
	for i := 0; i < len(t.ClaimTypes); i++ {
		if t.ClaimTypes[i].ID == lawSuit.ClaimTypeID {
			claimType = t.ClaimTypes[i].Name
			break
		}
	}
	// Статус - для фильтрации
	result["Lawsuit_ClaimTypeID"] = lawSuit.ClaimTypeID
	// Статус - для вывода в таблицу
	result["Lawsuit_ClaimType"] = claimType

	// TODO View ClaimTypes Улучшить поиск
	branch := ""
	for i := 0; i < len(c.Branches); i++ {
		if c.Branches[i].ID == lawSuit.BranchID {
			branch = c.Branches[i].Name
			break
		}
	}
	// Отделение - для фильтрации
	result["Lawsuit_BranchID"] = lawSuit.BranchID
	// Отделение - для вывода в таблицу
	result["Lawsuit_Branch"] = branch

	dbtSumFull := bc.Lawsuit.DebtSum +
		bc.Lawsuit.Penny +
		bc.Lawsuit.Penalty +
		bc.Lawsuit.StateDuty

	// Сумма процентов по 395 ГК РФ (руб.)
	if useFormat {
		result["Lawsuit_Percent_395"] = Currency(bc.Lawsuit.Percent395)
	} else {
		result["Lawsuit_Percent_395"] = bc.Lawsuit.Percent395
	}
	// Сумма процентов по 317.1 ГК РФ (руб.)
	if useFormat {
		result["Lawsuit_Percent_317"] = Currency(bc.Lawsuit.Percent317)
	} else {
		result["Lawsuit_Percent_317"] = bc.Lawsuit.Percent317
	}
	// Сумма договорной/законной неустойки (руб.)
	if useFormat {
		result["Lawsuit_Penalty"] = Currency(bc.Lawsuit.Penalty)
	} else {
		result["Lawsuit_Penalty"] = bc.Lawsuit.Penalty
	}
	// Пени по день фактической оплаты долга (руб.)
	if useFormat {
		result["Lawsuit_Penny"] = Currency(bc.Lawsuit.Penny)
	} else {
		result["Lawsuit_Penny"] = bc.Lawsuit.Penny
	}
	// Сумма госпошлины (руб.)
	if useFormat {
		result["Lawsuit_StateDuty"] = Currency(bc.Lawsuit.StateDuty)
	} else {
		result["Lawsuit_StateDuty"] = bc.Lawsuit.StateDuty
	}
	// Поступило денежных средств
	if useFormat {
		result["Lawsuit_ReceivedFunds"] = Currency(lawSuit.PaySum)
	} else {
		result["Lawsuit_ReceivedFunds"] = lawSuit.PaySum
	}
	// TODO Поле "Общий долг": Полная сумма долга
	if useFormat {
		result["Lawsuit_TotalDebt"] = Currency(dbtSumFull)
	} else {
		result["Lawsuit_TotalDebt"] = dbtSumFull
	}
	// TODO Поле "Основной долг": По счёт фактурам
	if useFormat {
		result["Lawsuit_MainDebt"] = Currency(lawSuit.InvoiceSum)
	} else {
		result["Lawsuit_MainDebt"] = lawSuit.InvoiceSum
	}
	// TODO Поле "Остаток долга": "Основной долг" - Поступило денежных средств
	if useFormat {
		result["Lawsuit_Balance"] = Currency(lawSuit.DebtSum)
	} else {
		result["Lawsuit_Balance"] = lawSuit.DebtSum
	}
	// TODO Колонка уведомление
	if useFormat {
		result["Lawsuit_Claim"] = Currency(lawSuit.InvoiceSum - lawSuit.PaySum)
	} else {
		result["Lawsuit_Claim"] = lawSuit.InvoiceSum - lawSuit.PaySum
	}
	// TODO Колонка претензия
	if useFormat {
		result["Lawsuit_Pretrial"] = "-" // Currency(lawSuit.InvoiceSum)
	} else {
		result["Lawsuit_Pretrial"] = "-" // lawSuit.InvoiceSum
	}

	if lawSuit.UnknownPayments {
		result["Lawsuit_UnknownPayments"] = 1
	} else {
		result["Lawsuit_UnknownPayments"] = 0
	}

	changes := make([]interface{}, 0)
	for i := 0; i < len(bc.ChangeItems); i++ {
		tmp := map[string]interface{}{
			"ID":        bc.ChangeItems[i].ID,
			"CreatedAt": formatTime(bc.ChangeItems[i].CreatedAt),
			"Action":    bc.ChangeItems[i].Key,
			"NewValue":  bc.ChangeItems[i].Value,
			"PrevValue": bc.ChangeItems[i].Prev,
		}
		changes = append(changes, tmp)
	}
	result["Lawsuit_Changes"] = changes

	contract := &lawSuit.Contract
	// ID
	result["Contract_ID"] = contract.ID
	// № Договор
	result["Contract_Number"] = contract.Number
	// Дата договора
	if useFormat {
		result["Contract_SignAt"] = formatDate(contract.SignAt)
	} else {
		result["Contract_SignAt"] = contract.SignAt
	}
	// Категория договора
	category := "Неизвестно"
	for i := 0; i < len(t.ContractCategoryTypes); i++ {
		if t.ContractCategoryTypes[i].ID == contract.CategoryID {
			category = t.ContractCategoryTypes[i].Name
			break
		}
	}
	result["Contract_CategoryID"] = contract.CategoryID
	result["Contract_Category"] = category
	// Статус
	if contract.Status == "" {
		result["Contract_Status"] = "Активен"
	} else {
		result["Contract_Status"] = contract.Status
	}
	// Почтовый адрес
	if contract.PostAddress == "" {
		result["Contract_PostAddress"] = "Не указан"
	} else {
		result["Contract_PostAddress"] = contract.PostAddress
	}
	// E-mail
	if contract.Email == "" {
		result["Contract_Email"] = "Не указан"
	} else {
		result["Contract_Email"] = contract.Email
	}
	// Валидность E-mail
	if contract.IsValidEmail {
		result["Contract_EmailValid"] = 1
	} else {
		result["Contract_EmailValid"] = 0
	}

	// Судебный отдел
	result["Contract_CuratorLegal"] = fmt.Sprintf("%v %v", contract.CuratorLegal.SecondName, contract.CuratorLegal.Name)
	result["Contract_CuratorLegal_Email"] = contract.CuratorLegal.Email
	result["Contract_CuratorLegal_Phone"] = contract.CuratorLegal.Phone
	// Расчётный отдел
	result["Contract_CuratorPayment"] = fmt.Sprintf("%v %v", contract.CuratorPayment.SecondName, contract.CuratorPayment.Name)
	result["Contract_CuratorPayment_Email"] = contract.CuratorPayment.Email
	result["Contract_CuratorPayment_Phone"] = contract.CuratorPayment.Phone
	// Договорной отдел
	result["Contract_CuratorContract"] = fmt.Sprintf("%v %v", contract.CuratorContract.SecondName, contract.CuratorContract.Name)
	result["Contract_CuratorContract_Email"] = contract.CuratorContract.Email
	result["Contract_CuratorContract_Phone"] = contract.CuratorContract.Phone
	// Отдел тех. аудита
	result["Contract_CuratorTechAudit"] = fmt.Sprintf("%v %v", contract.CuratorTechAudit.SecondName, contract.CuratorTechAudit.Name)
	result["Contract_CuratorTechAudit_Email"] = contract.CuratorTechAudit.Email
	result["Contract_CuratorTechAudit_Phone"] = contract.CuratorTechAudit.Phone
	// Куратор претензии
	result["Contract_CuratorClaim"] = fmt.Sprintf("%v %v", contract.CuratorClaim.SecondName, contract.CuratorClaim.Name)
	result["Contract_CuratorClaim_Email"] = contract.CuratorClaim.Email
	result["Contract_CuratorClaim_Phone"] = contract.CuratorClaim.Phone

	result["Contract_DaysToResolveClaim"] = contract.DaysToResolveClaim

	result["Contract_PaymentDay"] = 18
	for i := 0; i < len(contract.PaymentDays); i++ {
		if time.Now().After(contract.PaymentDays[i].DateFrom) &&
			time.Now().Before(contract.PaymentDays[i].DateTo) {
			result["Contract_PaymentDay"] = contract.PaymentDays[i].Day
			break
		}
	}

	paymentSchedules := make([]interface{}, 0)
	for i := 0; i < len(contract.PaymentSchedules); i++ {
		if contract.PaymentSchedules[i].ContractID == contract.ID {
			if useFormat {
				tmp := map[string]interface{}{
					"Day":     fmt.Sprintf("%d число", contract.PaymentSchedules[i].Day),
					"Percent": fmt.Sprintf("%d %%", contract.PaymentSchedules[i].Percent),
				}
				paymentSchedules = append(paymentSchedules, tmp)
			} else {
				tmp := map[string]interface{}{
					"Day":     fmt.Sprintf("%d", contract.PaymentSchedules[i].Day),
					"Percent": fmt.Sprintf("%d", contract.PaymentSchedules[i].Percent),
				}
				paymentSchedules = append(paymentSchedules, tmp)
			}
		}
	}
	if len(paymentSchedules) == 0 {
		tmp := map[string]interface{}{
			"Day":     fmt.Sprintf("18"),
			"Percent": fmt.Sprintf("100 %%"),
		}
		paymentSchedules = append(paymentSchedules, tmp)
	}
	// Срок оплаты по договору
	result["Contract_PaymentSchedules"] = paymentSchedules

	invoices := make([]interface{}, 0)
	totalSum := 0.0
	totalDebtSum := 0.0
	totalPayment := 0.0
	totalBalance := 0.0
	for i := 0; i < len(bc.Invoices); i++ {
		// TODO Пока аналитики вообще скрываю
		if strings.Trim(bc.Invoices[i].Document.Analytics, " ") != "" {
			continue
		}

		payment := 0.0
		for j := 0; j < len(bc.Payments); j++ {
			if bc.Invoices[i].ID == bc.Payments[j].InvoiceID {
				payment += bc.Payments[j].Sum
			}
		}

		note := bc.Invoices[i].Document.Note
		if note == "" {
			note = "Не задано"
		}

		tmp := map[string]interface{}{
			"ID":          bc.Invoices[i].ID,
			"ClaimNumber": lawSuit.NumberClaim,                            // Поле "Претензия"
			"Date":        formatDate(bc.Invoices[i].Document.DocumentAt), // Поле "Дата С/Ф"
			"Number":      bc.Invoices[i].Document.Number,                 // Поле "Номер С/Ф"
			"Type":        bc.Invoices[i].Document.Analytics,              // Поле "Тип начисления"
			"Count":       bc.Invoices[i].Count,                           // Кол-во кВт
			"Sum":         Currency(bc.Invoices[i].Sum),                   // Поле "Начислено"
			"DebtSum":     Currency(bc.Invoices[i].Sum - payment),         // Поле "Долг в претензии"
			"Payment":     Currency(payment),                              // Поле "Оплачено"
			"Balance":     Currency(bc.Invoices[i].Sum - payment),         // Поле "Остаток"
			"Note":        note,                                           // Поле "Примечание"
		}

		// TODO Аналитики нужно считать отдельно
		if strings.Trim(bc.Invoices[i].Document.Analytics, " ") == "" {
			totalSum += bc.Invoices[i].Sum
			totalDebtSum += bc.Invoices[i].Sum - payment
			totalPayment += payment
			totalBalance += bc.Invoices[i].Sum - payment
		}

		invoices = append(invoices, tmp)
	}
	// Счета фактуры по данному договору
	result["Contract_Invoices"] = invoices
	// Суммы счетов фактур по данному договору
	result["Contract_TotalInvoices"] = map[string]interface{}{
		"Sum":     Currency(totalSum),     // Поле "Начислено"
		"DebtSum": Currency(totalDebtSum), // Поле "Долг в претензии"
		"Payment": Currency(totalPayment), // Поле "Оплачено"
		"Balance": Currency(totalBalance), // Поле "Остаток"
	}

	result["Lawsuit_Period"] = bc.Lawsuit.ClaimPeriodStr

	payments := make([]interface{}, 0)
	totalSum = 0.0
	totalDebtSum = 0.0
	totalPayment = 0.0
	totalBalance = 0.0
	totalUnknownPayment := 0.0
	for i := 0; i < len(bc.Payments); i++ {
		// TODO Пока аналитики вообще скрываю
		if strings.Trim(bc.Payments[i].Document.Analytics, " ") != "" {
			continue
		}

		invoice := 0.0
		for j := 0; j < len(bc.Invoices); j++ {
			if bc.Payments[i].InvoiceID == bc.Invoices[j].ID {
				invoice += bc.Invoices[j].Sum
				break
			}
		}

		note := bc.Payments[i].Document.Note
		if note == "" {
			note = "Не задано"
		}

		tmp := map[string]interface{}{
			"ID":          bc.Payments[i].ID,
			"InvoiceID":   bc.Payments[i].InvoiceID,                       // Ссылка на С/Ф
			"ClaimNumber": lawSuit.NumberClaim,                            // Поле "Претензия"
			"Date":        formatDate(bc.Payments[i].Document.DocumentAt), // Поле "Дата С/Ф"
			"Number":      bc.Payments[i].Document.Number,                 // Поле "Номер С/Ф"
			"Type":        bc.Payments[i].Document.Analytics,              // Поле "Тип начисления"
			"Sum":         Currency(invoice),                              // Поле "Начислено"
			"DebtSum":     Currency(invoice - bc.Payments[i].Sum),         // Поле "Долг в претензии"
			"Payment":     Currency(bc.Payments[i].Sum),                   // Поле "Оплачено"
			"Balance":     Currency(invoice - bc.Payments[i].Sum),         // Поле "Остаток"
			"Note":        note,                                           // Поле "Примечание"
		}

		// TODO Аналитики нужно считать отдельно
		if strings.Trim(bc.Payments[i].Document.Analytics, " ") == "" {
			//  Платежи без фактур
			if bc.Payments[i].InvoiceID == 0 {
				totalUnknownPayment += bc.Payments[i].Sum
			} else {
				totalSum += invoice
				totalDebtSum += invoice - bc.Payments[i].Sum
				totalPayment += bc.Payments[i].Sum
				totalBalance += invoice - bc.Payments[i].Sum
			}
		}

		payments = append(payments, tmp)
	}
	// Платёжные документы по данному договору
	result["Contract_Payments"] = payments
	// Суммы платёжных документов по данному договору
	result["Contract_TotalPayments"] = map[string]interface{}{
		"Sum":     Currency(totalSum),     // Поле "Начислено"
		"DebtSum": Currency(totalDebtSum), // Поле "Долг в претензии"
		"Payment": Currency(totalPayment), // Поле "Оплачено"
		"Balance": Currency(totalBalance), // Поле "Остаток"
	}
	// Суммы нераспознанных платёжных документов по данному договору
	result["Contract_TotalUnknownPayments"] = map[string]interface{}{
		"Payment": Currency(totalUnknownPayment), // Поле "Сумма"
	}

	organization := &contract.Organization
	// ID
	result["Organization_ID"] = organization.ID
	// Email ЮЛ
	if organization.Email == "" {
		result["Organization_Email"] = "Не указан"
	} else {
		result["Organization_Email"] = organization.Email
	}
	// Наименование ЮЛ
	result["Organization_Name"] = organization.Name
	// ИНН ЮЛ
	result["Organization_INN"] = organization.INN
	// КПП ЮЛ
	result["Organization_KPP"] = organization.KPP
	// Юридический адрес ЮЛ
	result["Organization_LegalAddress"] = organization.LegalAddress
	// Категория организации
	category = "Неизвестно"
	for i := 0; i < len(t.OrganizationCategoryTypes); i++ {
		if t.OrganizationCategoryTypes[i].ID == organization.CategoryID {
			category = t.OrganizationCategoryTypes[i].Name
			break
		}
	}
	result["Organization_CategoryID"] = organization.CategoryID
	result["Organization_Category"] = category
	// Ликвидность организации
	if organization.IsLiquidated {
		result["Organization_Liquidity"] = "Ликвидирован"
	} else {
		result["Organization_Liquidity"] = "Действующий"
	}
	if organization.IsBankrupt {
		result["Organization_Bankrupt"] = "Банкрот"
	} else {
		result["Organization_Bankrupt"] = "Действующий"
	}

	FileMail := ""
	FileMailName := ""
	FileClaim := ""
	FileClaimName := ""
	FileClaimDetail := ""
	FileClaimDetailName := ""
	for i := 0; i < len(bc.Files); i++ {
		if strings.Contains(bc.Files[i].Name, "Письмо") {
			FileMail = bc.Files[i].FileID
			FileMailName = bc.Files[i].FileName
		}

		if strings.Contains(bc.Files[i].Name, "Претензия") {
			FileClaim = bc.Files[i].FileID
			FileClaimName = bc.Files[i].FileName
		}

		if strings.Contains(bc.Files[i].Name, "Реестр") {
			FileClaimDetail = bc.Files[i].FileID
			FileClaimDetailName = bc.Files[i].FileName
		}
	}
	result["File_Mail"] = FileMail
	result["File_MailName"] = FileMailName
	result["File_Claim"] = FileClaim
	result["File_ClaimName"] = FileClaimName
	result["File_ClaimDetail"] = FileClaimDetail
	result["File_ClaimDetailName"] = FileClaimDetailName

	// TODO Переделать под нормальные статусы
	if lawSuit.NotifyClaimDone {
		result["Notify_StatusDebt"] = fmt.Sprintf("%v %q %v", 1401, contract.Email, "Доставлено успешно")
	} else if lawSuit.NotifyClaimAt.Before(time.Now().AddDate(-10, 1, 1)) {
		result["Notify_StatusDebt"] = fmt.Sprintf("%v %q %v", 1401, contract.Email, "Ожидание")
	} else if !contract.IsValidEmail {
		result["Notify_StatusDebt"] = fmt.Sprintf("%v %q %v", 1401, contract.Email, "Не доставлено (недоступен канал)")
	} else if contract.Email == "" {
		result["Notify_StatusDebt"] = fmt.Sprintf("%v %q %v", 1401, "email не задан", "Не доставлено (отсутствует канал)")
	}

	// TODO Переделать под нормальные статусы
	if lawSuit.NotifyPretrialDone {
		result["Notify_StatusPretrial"] = fmt.Sprintf("%v %q %v", 1401, contract.Email, "Доставлено успешно")
	} else if lawSuit.NotifyPretrialAt.Before(time.Now().AddDate(-10, 1, 1)) {
		result["Notify_StatusPretrial"] = fmt.Sprintf("%v %q %v", 1401, contract.Email, "Ожидание")
	} else if !contract.IsValidEmail {
		result["Notify_StatusPretrial"] = fmt.Sprintf("%v %q %v", 1401, contract.Email, "Не доставлено (недоступен канал)")
	} else if contract.Email == "" {
		result["Notify_StatusPretrial"] = fmt.Sprintf("%v %q %v", 1401, "email не задан", "Не доставлено (отсутствует канал)")
	}

	result["Notify_ClaimChannel"] = lawSuit.NotifyClaimChannel
	result["Notify_ClaimCode"] = lawSuit.NotifyClaimCode
	result["Notify_ClaimDone"] = lawSuit.NotifyClaimDone
	result["Notify_ClaimMailingCode"] = lawSuit.NotifyClaimMailingCode
	result["Notify_PretrialChannel"] = lawSuit.NotifyPretrialChannel
	result["Notify_PretrialCode"] = lawSuit.NotifyPretrialCode
	result["Notify_PretrialDone"] = lawSuit.NotifyPretrialDone
	result["Notify_PretrialMailingCode"] = lawSuit.NotifyPretrialMailingCode

	if useFormat {
		result["Notify_ClaimAt"] = formatDate(lawSuit.NotifyClaimAt)
		result["Notify_PretrialAt"] = formatDate(lawSuit.NotifyPretrialAt)
	} else {
		result["Notify_ClaimAt"] = lawSuit.NotifyClaimAt
		result["Notify_PretrialAt"] = lawSuit.NotifyPretrialAt
	}

	return result, nil
}

// ClaimWorkView выборка
func ClaimWorkView(cw *ClaimWork, c *CommonRef, t *TypeRef, useFormat bool) (object_view.ViewMap, error) {
	result := make(object_view.ViewMap, 0)

	if c == nil {
		return result, nil
	}
	if len(cw.BriefCases) == 0 {
		return result, nil
	}

	for i := 0; i < len(cw.BriefCases); i++ {
		v, err := BriefCaseView(&cw.BriefCases[i], c, t, useFormat)
		if err != nil {
			return result, fmt.Errorf("ClaimWorkView, BriefCaseView[%v], Error: %v", i, err)
		}
		result.Append(fmt.Sprintf("%v", i), v)
	}

	return result, nil
}

// MailTemplateView параметры шаблона для письма
func MailTemplateView(bc *BriefCase) (map[string]interface{}, error) {
	result := make(map[string]interface{})

	if bc == nil {
		return result, fmt.Errorf("MailTemplateView, BriefCase is nil")
	}

	dbtSumFull := bc.Lawsuit.DebtSum +
		bc.Lawsuit.Penny +
		bc.Lawsuit.Penalty +
		bc.Lawsuit.StateDuty

	subject := "Уведомление о задолженности"
	template := "CLAIMDEBT"
	channel := "1401" // электронная почта
	mailingCode := bc.Lawsuit.NotifyClaimMailingCode

	if bc.Lawsuit.Status.Code != "1" {
		subject = "Досудебная претензия"
		template = "CLAIMPRETRIAL"
		channel = "1401"
		// channel := "1406" // Почта России
		mailingCode = bc.Lawsuit.NotifyPretrialMailingCode
	}

	// Курск - https://lkul-kursk.atomsbt.ru/
	// Мурманск – https://lkul-murmansk.atomsbt.ru/
	// Тверь - https://lkul-tver.atomsbt.ru/
	// Смоленск - https://lkul-smolensk.atomsbt.ru/
	// Хакасия - https://lkul-khakasia.atomsbt.ru/
	lkLink := bc.Lawsuit.Branch.PersonalAreaLink

	// {"cmdType":"mailing","mailing":{"head":{"protVer":"1.0.1","sysId":"RAPIRA","sourceSystem":"BS012","created":"2023-01-23 13:03:59","minProcVer":"1.0.1","senderVer":"1.0.1"},"command":"init","mailingList":[{"mailingCode":"TEST5_CLAIMDEBT20230119T132742","mailingPhaseCode":"1/1","startAt":"","timezone":"0","endAt":""}]}}

	// {"cmdType":"customerTemplateMessage","customerTemplateMessage":{"head":{"protVer":"1.0.1","sysId":"RAPIRA","sourceSystem":"BS012","created":"2023-01-23 13:03:59","minProcVer":"1.0.1","senderVer":"1.0.1"},"mailingCode":"TEST5_CLAIMDEBT20230119T132742","mailingPhaseCode":"1/1","templateCode":"CLAIMPRETRIAL","channelCode":"1401","fieldList":[{"contactInfo":"nechaevaa@atomsbt.ru","userId":"","userAddress":"183039, г. Мурманск, ул. Новое Плато, д.5, кв.60","persAcc":"5140145126","isOrganisationAcc":"1","header":"Досудебная претензия","dbtDate":"2022-09-30","dbtSum":"12 372,12","infoPhone":"9021356077","claimDate":"25.12.2022","claimPretrialDate":"28.12.2022","contractDate":"15.11.2018","contractNumber":"5140145126","dbtSumPeriod":"12 372,12","dbtSumFull":"14 460,29","lkLink":"https://lkul-murmansk.atomsbt.ru/","organisation":"","attachment":""}]}}

	result["StageCode"] = bc.Lawsuit.Stage.Code                          // Этап
	result["StatusCode"] = bc.Lawsuit.Status.Code                        // Статус
	result["mailingCode"] = mailingCode                                  // TODO Код рассылка "CLAIM20221004T112001" Обязательно в таком формате
	result["templateCode"] = template                                    // TODO Код шаблона (имя-строка на стороне уведомлений)
	result["channelCode"] = channel                                      // TODO Канал доставки
	result["lkLink"] = lkLink                                            // Ссылка на Личный кабинет
	result["claimDate"] = formatDate(bc.Lawsuit.ClaimAt)                 // Дата формирования претензии+3 к.д.
	result["claimPretrialDate"] = formatDate(bc.Lawsuit.PretrialAt)      // Дата формирования досудебной претензии+5 к.д.
	result["contactInfo"] = bc.Lawsuit.Contract.Email                    // Endpoint абонента
	result["contractDate"] = formatDate(bc.Lawsuit.Contract.SignAt)      // Дата договора
	result["contractNumber"] = bc.Lawsuit.Contract.Number                // Номер договора
	result["dbtDate"] = bc.Lawsuit.DateFrom.Format("2006-01-02")         // Период
	result["dbtDateStr"] = russianDate(bc.Lawsuit.DateFrom, true)        // TODO Переделать на строку Период строкой
	result["dbtSum"] = Currency(bc.Lawsuit.DebtSum)                      // TODO Сумма образовавшейся задолженности
	result["dbtSumFull"] = Currency(dbtSumFull)                          // Общая сумма долга
	result["dbtSumPeriod"] = Currency(bc.Lawsuit.DebtSum)                // Сумма долга за период
	result["mailingSubject"] = subject                                   // Заголовок письма
	result["infoPhone"] = bc.Lawsuit.Contract.Organization.Phone         // Телефон абонента
	result["isOrganisationAcc"] = "1"                                    // true // Если организация
	result["organisation"] = bc.Lawsuit.Contract.Organization.Name       // Организация
	result["persAcc"] = bc.Lawsuit.Contract.Number                       // Лицевой счёт / номер договора
	result["userAddress"] = bc.Lawsuit.Contract.Organization.PostAddress // Почтовый адрес абонента

	attachments := make([]interface{}, 0)
	for i := 0; i < len(bc.Files); i++ {
		if !strings.Contains(bc.Files[i].Name, "Претензия") {
			continue
		}
		tmp := map[string]interface{}{
			"bucketName":    "claim",
			"fileName":      bc.Files[i].FileName,
			"fileExtension": bc.Files[i].Extension,
			"fileSizeByte":  bc.Files[i].Size,
			"pathToFile":    "",
			"eTag":          "",
		}
		attachments = append(attachments, tmp)
	}
	result["attachments"] = attachments // Вложения

	return result, nil
}

// WhiteListView Белый список договоров
func WhiteListView(c *CommonRef) ([]interface{}, error) {
	result := make([]interface{}, 0)

	for i := 0; i < len(c.WhiteList); i++ {
		item := c.WhiteList[i]

		tmp := map[string]interface{}{
			"ID":                   item.ID,
			"ContractNumber":       item.Contract.Number,
			"CreatedAt":            formatTime(item.CreatedAt),
			"CreatedBy":            item.CreatedBy.Name,
			"DateFrom":             formatTime(item.DateFrom),
			"DateTo":               formatTime(item.DateTo),
			"EDMSLink":             item.EDMSLink,
			"ModifiedAt":           formatTime(item.ModifiedAt),
			"ModifiedBy":           item.ModifiedBy.Name,
			"Note":                 item.Note,
			"OrganizationCategory": item.Contract.Organization.CategoryID,
			"OrganizationINN":      item.Contract.Organization.INN,
			"Reason":               item.Reason.Name,
		}

		result = append(result, tmp)
	}

	return result, nil
}

// ===========================================================================

func Currency(number float64) string {
	tmp := fmt.Sprintf("%d", int64(number))
	res := ""
	j := 0
	for i := len(tmp) - 1; i >= 0; i-- {
		j++
		res = string(tmp[i]) + res
		if j > 0 && j%3 == 0 {
			res = " " + res
		}
	}

	tmp1 := fmt.Sprintf("%.2f", number)
	tmp2 := strings.Split(tmp1, ".")
	res = res + "." + tmp2[1]

	// fmt.Println(res)
	res = strings.Trim(res, " ")
	return strings.ReplaceAll(res, ".", ",")
}

func formatDate(date time.Time) string {
	return date.Format("02.01.2006")
}

func formatTime(date time.Time) string {
	return date.Format("02.01.2006 15:04:05")
}

// PeriodDates Диапазон дат строкой
// Январь 2022 - Март 2022
// Январь 2022, Март 2022, Июнь 2022
// Январь 2022, Март 2022 - Июнь 2022
// TODO Реализовать сложный (составной) период
func PeriodDates(date1 time.Time, date2 time.Time) string {
	if date1.Year() == date2.Year() && date1.Month() == date2.Month() {
		return fmt.Sprintf("%v",
			date2.Format("01.2006"))
	} else if date1.Before(date2) {
		return fmt.Sprintf("%v-%v",
			date1.Format("01.2006"),
			date2.Format("01.2006"))
	} else {
		return fmt.Sprintf("%v-%v",
			date2.Format("01.2006"),
			date1.Format("01.2006"))
	}
}

func russianDate(date time.Time, long bool) string {
	months := make(map[int]string, 0)
	if long {
		months[1] = "январь"
		months[2] = "февраль"
		months[3] = "март"
		months[4] = "апрель"
		months[5] = "май"
		months[6] = "июнь"
		months[7] = "июль"
		months[8] = "август"
		months[9] = "сентябрь"
		months[10] = "октябрь"
		months[11] = "ноябрь"
		months[12] = "декабрь"
	} else {
		months[1] = "янв"
		months[2] = "февр"
		months[3] = "март"
		months[4] = "апр"
		months[5] = "май"
		months[6] = "июнь"
		months[7] = "июль"
		months[8] = "авг"
		months[9] = "сент"
		months[10] = "окт"
		months[11] = "нояб"
		months[12] = "дек"
	}

	m := date.Month()
	y := date.Year()

	return fmt.Sprintf("%v %v", months[int(m)], y)
}
