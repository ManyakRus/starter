// модуль для хранения постоянных переменных, констант
package constants

import (
	model "gitlab.aescorp.ru/dsp_dev/claim/common/object_model"
	"time"
)

var Loc = time.Local

// CONNECTION_ID - ИД в БД Рапира в таблице connections
var CONNECTION_ID int64 = 3 //7

// BRANCH_ID - ИД в БД Рапира в таблице branches
var BRANCH_ID int64 = 2 //20954

// CONNECTION - объект Соединение, настроенный
var CONNECTION = model.Connection{ID: CONNECTION_ID, BranchId: BRANCH_ID, IsLegal: true, Server: "10.1.9.153", Port: "5432", DbName: "kol_atom_ul_uni", DbScheme: "stack", Login: "", Password: ""}
