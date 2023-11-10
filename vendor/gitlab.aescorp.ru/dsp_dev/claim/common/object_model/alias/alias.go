// Package alias -- специальные типы РАПИРы
package alias

// PaymentId -- ID платёжки
type PaymentId int64

// InvoiceId -- ID счёт-фактуры
type InvoiceId int64

// LawsuitId -- ID претензии
type LawsuitId int64

// LawsuitNumber -- номер претензии
type LawsuitNumber string

// ClaimNumber -- Номер дела
type ClaimNumber string

// TrialNumber -- Номер иска
type TrialNumber string

// ContractNumber -- Номер договора
type ContractNumber string

// FrontDate -- специальный тип даты для фронта
type FrontDate string

// FrontTime -- специальный тип даты-времени для фронта
type FrontTime string

// PaymentRegisteredAt -- тип даты времени при регистрации в системе
type PaymentRegisteredAt string

// IsAfterNotify -- признак регистрации документа после уведомления
type IsAfterNotify bool
