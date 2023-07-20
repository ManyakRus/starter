package object_model

import (
	"fmt"
	"time"

	"github.com/vmihailenco/msgpack/v5"
	"gitlab.aescorp.ru/dsp_dev/claim/common/object_model/alias"
)

// Lawsuit Дело. Объединяет весь набор данных по конкретному должнику.
type Lawsuit struct {
	CommonStruct
	GroupStruct
	NameStruct
	Branch                    Branch              `json:"branch"                  gorm:"-:all"`
	BranchID                  int64               `json:"branch_id"               gorm:"column:branch_id;default:null"`
	Chance                    string              `json:"chance"                  gorm:"column:chance;default:\"\""`
	ClaimAt                   time.Time           `json:"claim_at"                gorm:"column:claim_at;default:null"` // Уведомление о задолженности. Оплатить до.
	ClaimPeriodStr            string              `json:"claim_period_str"        gorm:"column:claim_period_str;default:\"\""`
	ClaimType                 ClaimType           `json:"claim_type"              gorm:"-:all"` // Тип задолженности
	ClaimTypeID               int64               `json:"claim_type_id"           gorm:"column:claim_type_id;default:null"`
	ClosedAt                  time.Time           `json:"closed_at"               gorm:"column:closed_at;default:null"`
	Contract                  Contract            `json:"contract"                gorm:"-:all"` // Договор
	ContractID                int64               `json:"contract_id"             gorm:"column:contract_id;default:null"`
	ControlledAt              time.Time           `json:"controlled_at"           gorm:"column:controlled_at;default:null"`
	Court                     Court               `json:"court"                   gorm:"-:all"`
	CourtID                   int64               `json:"court_id"                gorm:"column:court_id;default:null"`
	CuratorClaim              Employee            `json:"curator_claim"           gorm:"-:all"`
	CuratorClaimID            int64               `json:"curator_claim_id"        gorm:"column:curator_claim_id;default:null"`
	CuratorContract           Employee            `json:"curator_contract"        gorm:"-:all"`
	CuratorContractID         int64               `json:"curator_contract_id"     gorm:"column:curator_contract_id;default:null"`
	CuratorLegal              Employee            `json:"curator_legal"           gorm:"-:all"`
	CuratorLegalID            int64               `json:"curator_legal_id"        gorm:"column:curator_legal_id;default:null"`
	CuratorPayment            Employee            `json:"curator_payment"         gorm:"-:all"`
	CuratorPaymentID          int64               `json:"curator_payment_id"      gorm:"column:curator_payment_id;default:null"`
	CuratorTechAudit          Employee            `json:"curator_tech_audit"      gorm:"-:all"`
	CuratorTechAuditID        int64               `json:"curator_tech_audit_id"   gorm:"column:curator_tech_audit_id;default:null"`
	DateFrom                  time.Time           `json:"date_from"               gorm:"column:date_from;default:null"`
	DateTo                    time.Time           `json:"date_to"                 gorm:"column:date_to;default:null"`
	DebtSum                   float64             `json:"debt_sum"                gorm:"column:debt_sum;default:0"`             // Общая сумма долга
	DebtSumSentClaim          float64             `json:"debt_sum_sent_claim"     gorm:"column:debt_sum_sent_claim;default:0"`  // Общий долг при отправке претензии, руб.
	DebtSumSentNotify         float64             `json:"debt_sum_sent_notify"    gorm:"column:debt_sum_sent_notify;default:0"` // Общий долг при отправке уведомления, руб.
	InvoiceSum                float64             `json:"invoice_sum"             gorm:"column:invoice_sum;default:0"`          // Сумма долга за период
	IsClosed                  bool                `json:"is_closed"               gorm:"column:is_closed;default:null"`
	NotifyClaimAt             time.Time           `json:"notify_claim_at"         gorm:"column:notify_claim_at;default:null"`                   // Уведомление о задолженности. Дата отправки.
	NotifyClaimChannel        int                 `json:"notify_claim_channel"    gorm:"column:notify_claim_channel;default:null"`              // Уведомление о задолженности. Канал отправки.
	NotifyClaimCode           int                 `json:"notify_claim_code"       gorm:"column:notify_claim_code;default:null"`                 // Уведомление о задолженности. Код доставки из НСИ.
	NotifyClaimDone           bool                `json:"notify_claim_done"       gorm:"column:notify_claim_done;default:null"`                 // Уведомление о задолженности. Факт отправки.
	NotifyClaimMailingCode    string              `json:"notify_claim_mailing_code" gorm:"column:notify_claim_mailing_code;default:null"`       // Уведомление о задолженности. Уникальный код отправки.
	NotifyPretrialAt          time.Time           `json:"notify_pretrial_at"      gorm:"column:notify_pretrial_at;default:null"`                // Досудебная претензия. Дата отправки.
	NotifyPretrialChannel     int                 `json:"notify_pretrial_channel" gorm:"column:notify_pretrial_channel;default:null"`           // Досудебная претензия. Канал отправки.
	NotifyPretrialCode        int                 `json:"notify_pretrial_code"    gorm:"column:notify_pretrial_code;default:null"`              // Досудебная претензия. Код доставки из НСИ.
	NotifyPretrialDone        bool                `json:"notify_pretrial_done"    gorm:"column:notify_pretrial_done;default:null"`              // Досудебная претензия. Факт отправки.
	NotifyPretrialMailingCode string              `json:"notify_pretrial_mailing_code" gorm:"column:notify_pretrial_mailing_code;default:null"` // Досудебная претензия. Уникальный код отправки.
	Number                    alias.LawsuitNumber `json:"number"                  gorm:"column:number;default:\"\""`
	NumberClaim               alias.ClaimNumber   `json:"number_claim"            gorm:"column:number_claim;default:\"\""`
	NumberTrial               string              `json:"number_trial"            gorm:"column:number_trial;default:\"\""`
	PaySum                    float64             `json:"pay_sum"                 gorm:"column:pay_sum;default:0"` // Платежи
	Penalty                   float64             `json:"penalty"                 gorm:"column:penalty;default:0"`
	Penny                     float64             `json:"penny"                   gorm:"column:penny;default:0"`
	Percent317                float64             `json:"percent_317"             gorm:"column:percent_317;default:0"`
	Percent395                float64             `json:"percent_395"             gorm:"column:percent_395;default:0"`
	PretrialAt                time.Time           `json:"pretrial_at"             gorm:"column:pretrial_at;default:null"` // Досудебная претензия. Оплатить до.
	ProcessStartedAt          time.Time           `json:"process_started_at"      gorm:"column:process_started_at;default:null"`
	ProcessKey                string              `json:"process_key"             gorm:"column:process_key;default:\"\""`
	Reason                    LawsuitReasonType   `json:"reason"                  gorm:"-:all"`
	ReasonID                  int64               `json:"reason_id"               gorm:"column:reason_id;default:null"`
	RestrictSum               float64             `json:"restrict_sum"            gorm:"column:restrict_sum;default:0"`
	Stage                     LawsuitStageType    `json:"stage"                   gorm:"-:all"` // Этап
	StageAt                   time.Time           `json:"stage_at"                gorm:"column:stage_at;default:null"`
	StageID                   int64               `json:"stage_id"                gorm:"column:stage_id;default:null"`
	StateDuty                 float64             `json:"state_duty"              gorm:"column:state_duty;default:0"` // Пошлина
	Status                    LawsuitStatusType   `json:"status"                  gorm:"-:all"`                       // Статус
	StatusAt                  time.Time           `json:"status_at"               gorm:"column:status_at;default:null"`
	StatusID                  int64               `json:"status_id"               gorm:"column:status_id;default:null"`
	Tag                       string              `json:"tag"                     gorm:"column:tag;default:\"\""`
	UnknownPayments           bool                `json:"unknown_payments"        gorm:"unknown_payments:tag;default:false"` // "С не разнесёнными платежами"
}

// ClaimNumber -- возвращает номер портфеля
func (sf *Lawsuit) ClaimNumber() alias.ClaimNumber {
	return sf.NumberClaim
}

// NewLawsuit Новый объект дела
func NewLawsuit(contractNumber alias.ContractNumber) Lawsuit {
	strClaimNumber := fmt.Sprintf("ПР_%s_%s", time.Now().Format("200601-02"), contractNumber)
	strContractNumber := fmt.Sprintf("ПФ_%s_%s", time.Now().Format("200601-02"), contractNumber)
	sf := Lawsuit{
		Number:      alias.LawsuitNumber(strContractNumber),
		NumberClaim: alias.ClaimNumber(strClaimNumber),
		NumberTrial: fmt.Sprintf("ПИ_%s_%s", time.Now().Format("200601-02"), contractNumber),
	}
	return sf
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
