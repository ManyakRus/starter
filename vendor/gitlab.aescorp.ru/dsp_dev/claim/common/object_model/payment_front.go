package object_model

// PaymentFront -- платежка для отдачи на фронт
type PaymentFront struct {
	Id_                int64  `json:"ID"`                       // ID платежа
	InvoiceId_         int64  `json:"InvoiceID"`                // ID с/ф
	ClaimNumber_       string `json:"ClaimNumber"`              // Номер претензии
	DatePayAt_         string `json:"Date"`                     // Дата оплаты п/п
	DateRegistreAt_    string `json:"DistributionDate"`         // Дата попадания платежа в систему
	IsAfterNotify_     bool   `json:"is_payment_after_created"` // Если платеж был после уведомления
	Number_            string `json:"Number"`                   // Номер платежа
	Type_              string `json:"Type"`                     // Тип документа
	InvoiceSum_        string `json:"Sum"`                      // Сумма (должна быть пустой для платежа)
	InvoiceCorrectSum_ string `json:"Correction"`               // Сумма корректировки если п/п оказалось с/ф
	InvoiceDebtSum_    string `json:"DebtSum"`                  // Сумма задолженности если п/п оказалось с/ф (должно быть пустым для платежа)
	Sum_               string `json:"Payment"`                  // Сумма платежа
	InvoiceBalance_    string `json:"Balance"`                  // Остаток по с/ф после всех оплат (должна быть пустой для платежа)
	Note_              string `json:"Note"`                     // Комментарий к платежу
}
