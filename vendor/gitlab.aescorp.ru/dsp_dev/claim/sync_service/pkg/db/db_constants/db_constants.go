package db_constants

import (
	"errors"
)

const CONNECTION_ID_TEST = 3

const TIMEOUT_DB_SECONDS = 30

const TEXT_RECORD_NOT_FOUND = "record not found"

const TextCrudIsNotInit = "Need initializate crud with InitCrudTransport_DB() function at first."

var ErrorCrudIsNotInit error = errors.New(TextCrudIsNotInit)
