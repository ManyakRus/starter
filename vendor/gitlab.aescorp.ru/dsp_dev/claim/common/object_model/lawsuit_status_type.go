package object_model

import (
	"encoding/json"
	"github.com/vmihailenco/msgpack/v5"
	"reflect"
)

// versionLawsuitStatusType - версия структуры модели, с учётом имен и типов полей
var versionLawsuitStatusType uint32

// crud_LawsuitStatusType - объект контроллер crud операций
var crud_LawsuitStatusType ICrud_LawsuitStatusType

// LawsuitStatusType Статусы дел (справочник).
type LawsuitStatusType struct {
	CommonStruct
	NameStruct
	Code string `json:"code" gorm:"column:code;default:0"`
}

type ICrud_LawsuitStatusType interface {
	Read(*LawsuitStatusType) error
	Save(*LawsuitStatusType) error
	Update(*LawsuitStatusType) error
	Create(*LawsuitStatusType) error
	Delete(*LawsuitStatusType) error
	Restore(*LawsuitStatusType) error
}

// TableName - возвращает имя таблицы в БД, нужен для gorm
func (m LawsuitStatusType) TableNameDB() string {
	return "lawsuit_status_types"
}

// NewLawsuitStatusType - возвращает новый	объект
func NewLawsuitStatusType() LawsuitStatusType {
	return LawsuitStatusType{}
}

// AsLawsuitStatusType - создаёт объект из упакованного объекта в массиве байтов
func AsLawsuitStatusType(b []byte) (LawsuitStatusType, error) {
	c := NewLawsuitStatusType()
	err := msgpack.Unmarshal(b, &c)
	if err != nil {
		return NewLawsuitStatusType(), err
	}
	return c, nil
}

// LawsuitStatusTypeAsBytes - упаковывает объект в массив байтов
func LawsuitStatusTypeAsBytes(m *LawsuitStatusType) ([]byte, error) {
	b, err := msgpack.Marshal(m)
	if err != nil {
		return nil, err
	}
	return b, nil
}

// GetStructVersion - возвращает версию модели
func (m LawsuitStatusType) GetStructVersion() uint32 {
	if versionLawsuitStatusType == 0 {
		versionLawsuitStatusType = CalcStructVersion(reflect.TypeOf(m))
	}

	return versionLawsuitStatusType
}

// GetModelFromJSON - создаёт модель из строки json
func (m *LawsuitStatusType) GetModelFromJSON(sModel string) error {
	var err error

	var bytes []byte
	bytes = []byte(sModel)

	err = json.Unmarshal(bytes, m)

	return err
}

// GetJSON - возвращает строку json из модели
func (m LawsuitStatusType) GetJSON() (string, error) {
	var ReturnVar string
	var err error

	bytes, err := json.Marshal(m)
	if err != nil {
		return ReturnVar, err
	}
	ReturnVar = string(bytes)
	return ReturnVar, err
}

//---------------------------- CRUD операции ------------------------------------------------------------

// Read - находит запись в БД по ID, и заполняет в объект
func (m *LawsuitStatusType) Read() error {
	err := crud_LawsuitStatusType.Read(m)

	return err
}

// Save - записывает объект в БД по ID
func (m *LawsuitStatusType) Save() error {
	err := crud_LawsuitStatusType.Save(m)

	return err
}

// Update - обновляет объект в БД по ID
func (m *LawsuitStatusType) Update() error {
	err := crud_LawsuitStatusType.Update(m)

	return err
}

// Create - создаёт объект в БД с новым ID
func (m *LawsuitStatusType) Create() error {
	err := crud_LawsuitStatusType.Create(m)

	return err
}

// Delete - устанавливает признак пометки удаления в БД
func (m *LawsuitStatusType) Delete() error {
	err := crud_LawsuitStatusType.Delete(m)

	return err
}

// Restore - снимает признак пометки удаления в БД
func (m *LawsuitStatusType) Restore() error {
	err := crud_LawsuitStatusType.Restore(m)

	return err
}

// SetCrudInterface - заполняет интерфейс crud: DB, GRPC, NRPC
func (m LawsuitStatusType) SetCrudInterface(crud ICrud_LawsuitStatusType) {
	crud_LawsuitStatusType = crud

	return
}

//---------------------------- конец CRUD операции ------------------------------------------------------------
