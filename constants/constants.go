// модуль для хранения постоянных переменных, констант
package constants

import (
	"time"
)

// LayoutDateTimeRus - формат текстовой даты и времени для России
var LayoutDateTimeRus = "02.01.2006 15:04:05"

// LayoutDateRus - формат текстовой даты для России
var LayoutDateRus = "02.01.2006"

var Loc = time.Local

// CONNECTION_ID - ИД в БД Рапира в таблице connections
var CONNECTION_ID int64 = 3 //7

// BRANCH_ID - ИД в БД Рапира в таблице branches
var BRANCH_ID int64 = 2 //20954

var TIME_ZONE = "Europe/Moscow"
