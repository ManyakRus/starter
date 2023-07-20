package object_model

import (
	"github.com/vmihailenco/msgpack/v5"
)

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

	// Досудебная претензия, которая на стадии мониторинга с направлением по e-mail
	ClaimsChannel1401    int    `json:"claims_channel_1401"      gorm:"-:all"`
	ClaimsChannel1401Sum string `json:"claims_channel_1401_sum"  gorm:"-:all"`

	// Досудебная претензия, которая на стадии мониторинга с направлением заказным письмом
	ClaimsChannel1406    int    `json:"claims_channel_1406"      gorm:"-:all"`
	ClaimsChannel1406Sum string `json:"claims_channel_1406_sum"  gorm:"-:all"`

	// Претензии с платежами не разнесенными на счет-фактуры
	ClaimsWithUnknown    int    `json:"claims_with_unknown"      gorm:"-:all"`
	ClaimsWithUnknownSum string `json:"claims_with_unknown_sum"  gorm:"-:all"`
}

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
