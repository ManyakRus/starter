// модуль для хранения постоянных переменных, констант
package constants

import (
	"gitlab.aescorp.ru/dsp_dev/claim/sync_service/pkg/db/tables/table_connections"
	"gitlab.aescorp.ru/dsp_dev/claim/sync_service/pkg/object_model/entities/connections"
	"time"
)

var Loc = time.Local

// CONNECTION_ID - ИД в БД Рапира в таблице connections
var CONNECTION_ID int64 = 3 //7

// BRANCH_ID - ИД в БД Рапира в таблице branches
var BRANCH_ID int64 = 2 //20954

// CONNECTION - объект Соединение, настроенный
var CONNECTION = connections.Connection{Table_Connection: table_connections.Table_Connection{ID: CONNECTION_ID, BranchID: BRANCH_ID, IsLegal: true, Server: "10.1.9.153", Port: "5432", DbName: "kol_atom_ul_uni", DbScheme: "stack", Login: "", Password: ""}}

var TIME_ZONE = "Europe/Moscow"
