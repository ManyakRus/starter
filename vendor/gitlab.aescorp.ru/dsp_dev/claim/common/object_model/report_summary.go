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

	// На этапе уведомления
	ClaimsStatus3    int    `json:"claims_status_3"      gorm:"-:all"`
	ClaimsStatus3Sum string `json:"claims_status_3_sum"  gorm:"-:all"`

	// Завершено на этапе претензии в связи с оплатой с начала года
	ClaimsStatus6    int    `json:"claims_status_6"      gorm:"-:all"`
	ClaimsStatus6Sum string `json:"claims_status_6_sum"  gorm:"-:all"`

	// Ожидает передачи в исковое производство
	ClaimsStatus8    int    `json:"claims_status_8"      gorm:"-:all"`
	ClaimsStatus8Sum string `json:"claims_status_8_sum"  gorm:"-:all"`

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
