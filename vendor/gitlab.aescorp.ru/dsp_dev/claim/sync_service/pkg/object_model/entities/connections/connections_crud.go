//Файл создан автоматически кодогенератором crud_generator
//Не изменяйте ничего здесь.

package connections

import (
	"encoding/json"
	"github.com/vmihailenco/msgpack/v5"
	"gitlab.aescorp.ru/dsp_dev/claim/sync_service/pkg/db/calc_struct_version"
	"gitlab.aescorp.ru/dsp_dev/claim/sync_service/pkg/db/db_constants"
	"reflect"
)

// versionConnection - версия структуры модели, с учётом имен и типов полей
var versionConnection uint32

// Crud_Connection - объект контроллер crud операций
var Crud_Connection ICrud_Connection

// интерфейс стандартных CRUD операций, для использования в DB или GRPC или NRPC
type ICrud_Connection interface {
	Read(*Connection) error
	Save(*Connection) error
	Update(*Connection) error
	Create(*Connection) error
	ReadFromCache(ID int64) (Connection, error)
	UpdateManyFields(*Connection, []string) error
	Update_BranchID(*Connection) error
	Update_DbName(*Connection) error
	Update_DbScheme(*Connection) error
	Update_IsLegal(*Connection) error
	Update_Login(*Connection) error
	Update_Name(*Connection) error
	Update_Password(*Connection) error
	Update_Port(*Connection) error
	Update_Server(*Connection) error
}

// TableName - возвращает имя таблицы в БД
func (m Connection) TableNameDB() string {
	return "connections"
}

// NewConnection - возвращает новый	объект
func NewConnection() Connection {
	return Connection{}
}

// AsConnection - создаёт объект из упакованного объекта в массиве байтов
func AsConnection(b []byte) (Connection, error) {
	c := NewConnection()
	err := msgpack.Unmarshal(b, &c)
	if err != nil {
		return NewConnection(), err
	}
	return c, nil
}

// ConnectionAsBytes - упаковывает объект в массив байтов
func ConnectionAsBytes(m *Connection) ([]byte, error) {
	b, err := msgpack.Marshal(m)
	if err != nil {
		return nil, err
	}
	return b, nil
}

// GetStructVersion - возвращает версию модели
func (m Connection) GetStructVersion() uint32 {
	if versionConnection == 0 {
		versionConnection = calc_struct_version.CalcStructVersion(reflect.TypeOf(m))
	}

	return versionConnection
}

// GetModelFromJSON - создаёт модель из строки json
func (m *Connection) GetModelFromJSON(sModel string) error {
	var err error

	var bytes []byte
	bytes = []byte(sModel)

	err = json.Unmarshal(bytes, m)

	return err
}

// GetJSON - возвращает строку json из модели
func (m Connection) GetJSON() (string, error) {
	var ReturnVar string
	var err error

	bytes, err := json.Marshal(m)
	if err != nil {
		return ReturnVar, err
	}
	ReturnVar = string(bytes)
	return ReturnVar, err
}

// ---------------------------- CRUD операции ------------------------------------------------------------

// Read - находит запись в БД по ID, и заполняет в объект
func (m *Connection) Read() error {
	if Crud_Connection == nil {
		return db_constants.ErrorCrudIsNotInit
	}

	err := Crud_Connection.Read(m)

	return err
}

// Save - записывает объект в БД по ID
func (m *Connection) Save() error {
	if Crud_Connection == nil {
		return db_constants.ErrorCrudIsNotInit
	}

	err := Crud_Connection.Save(m)

	return err
}

// Update - обновляет объект в БД по ID
func (m *Connection) Update() error {
	if Crud_Connection == nil {
		return db_constants.ErrorCrudIsNotInit
	}

	err := Crud_Connection.Update(m)

	return err
}

// Create - создаёт объект в БД с новым ID
func (m *Connection) Create() error {
	if Crud_Connection == nil {
		return db_constants.ErrorCrudIsNotInit
	}

	err := Crud_Connection.Create(m)

	return err
}

// ReadFromCache - находит запись в кэше или в БД по ID, и заполняет в объект
func (m *Connection) ReadFromCache(ID int64) (Connection, error) {
	Otvet := Connection{}
	var err error

	if Crud_Connection == nil {
		return Otvet, db_constants.ErrorCrudIsNotInit
	}

	Otvet, err = Crud_Connection.ReadFromCache(ID)

	return Otvet, err
}

// SetCrudInterface - заполняет интерфейс crud: DB, GRPC, NRPC
func (m Connection) SetCrudInterface(crud ICrud_Connection) {
	Crud_Connection = crud

	return
}

// UpdateManyFields - находит запись в БД по ID, и изменяет только нужные колонки
func (m *Connection) UpdateManyFields(MassNeedUpdateFields []string) error {
	if Crud_Connection == nil {
		return db_constants.ErrorCrudIsNotInit
	}

	err := Crud_Connection.UpdateManyFields(m, MassNeedUpdateFields)

	return err
}

// ---------------------------- конец CRUD операции ------------------------------------------------------------
