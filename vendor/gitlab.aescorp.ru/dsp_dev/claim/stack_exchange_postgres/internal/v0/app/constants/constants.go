// модуль для хранения постоянных переменных, констант
package constants

import (
	"time"

	"gitlab.aescorp.ru/dsp_dev/claim/sync_service/pkg/object_model/entities/connections"
)

// Layout - формат текстовой даты для загрузки из json
var Layout = "2006-01-02 15:04:05.999999999 Z0700 MST"

// Loc - Системная локальная time zone
var Loc = time.Local

// LocationUTC - time zone UTC
var LocationUTC, _ = time.LoadLocation("UTC")

// Date1Test - дата начала загрузки из СТЕК при запуске тестов
var Date1Test = time.Date(2023, time.September, 01, 0, 0, 0, 0, LocationUTC)

// Date2Test - дата окончания загрузки из СТЕК при запуске тестов
var Date2Test = time.Date(2023, time.September, 22, 23, 59, 59, 0, LocationUTC)

//var Date2Test = carbon.Time2Carbon(Date1Test).EndOfMonth().EndOfDay().Carbon2Time()

// SERVICE_NAME - название данного сервиса
var SERVICE_NAME = "stack_exchange"

// CONNECTION_ID - ИД в БД Рапира в таблице connections
var CONNECTION_ID int64 = 3 //7

// BRANCH_ID - ИД в БД Рапира в таблице branches
var BRANCH_ID int64 = 2 //20954

// NEED_UPGRADE_TEST - надо ли обновлять записи в БД Рапира, при запуске тестов
const NEED_UPGRADE_TEST = true

// CONNECTION - объект Соединение, настроенный
var CONNECTION = connections.Connection{ID: CONNECTION_ID, BranchID: BRANCH_ID, IsLegal: true}

// CAMUNDA

// CAMUNDA_JOBTYPE - имя задачи в CAMUNDA
var CAMUNDA_JOBTYPE = "stack_exchange"

// CAMUNDA_ID - ИД сервиса в CAMUNDA
var CAMUNDA_ID = "BS012_100_006"

// CAMUNDA_BPMNFILE - путь к файлу бизнес процесса .bpmn
// не нужен
var CAMUNDA_BPMNFILE = ""

//var CAMUNDA_BPMNFILE = "bpmn/claim_process.bpmn"
