package constants

import (
	"errors"
)

// CONNECTION_ID_TEST - Connection ID для тестов
const CONNECTION_ID_TEST = 3

// TIMEOUT_DB_SECONDS - время ожидания в секундах
const TIMEOUT_DB_SECONDS = 30

// TEXT_RECORD_NOT_FOUND - текст ошибки, если нет записи
const TEXT_RECORD_NOT_FOUND = "record not found"

// TextCrudIsNotInit - текст ошибки, если не инициализирован crud
const TextCrudIsNotInit = "Need initializate crud with InitCrudTransport_NRPC() function at first."

// ErrorCrudIsNotInit - ошибка, если не инициализирован crud
var ErrorCrudIsNotInit error = errors.New(TextCrudIsNotInit)
