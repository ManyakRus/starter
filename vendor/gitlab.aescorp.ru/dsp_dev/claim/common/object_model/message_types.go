package object_model

import (
	"encoding/json"
	"github.com/vmihailenco/msgpack/v5"
	"reflect"
)

// versionMessageType - версия структуры модели, с учётом имен и типов полей
var versionMessageType uint32

// crud_MessageType - объект контроллер crud операций
var crud_MessageType ICrud_MessageType

// MessageType - Типы сообщений
type MessageType struct {
	CommonStruct
	NameStruct
	Code int `json:"code"        gorm:"column:code;default:0"`
}

type ICrud_MessageType interface {
	Read(m *MessageType) error
	Save(m *MessageType) error
	Update(m *MessageType) error
	Create(m *MessageType) error
	Delete(m *MessageType) error
	Restore(m *MessageType) error
}

// TableName - возвращает имя таблицы в БД, нужен для gorm
func (m MessageType) TableNameDB() string {
	return "message_types"
}

// GetID - возвращает ID объекта
func (m MessageType) GetID() int64 {
	return m.ID
}

// NewMessageType - возвращает новый	объект
func NewMessageType() MessageType {
	return MessageType{}
}

// AsMessageType - создаёт объект из упакованного объекта в массиве байтов
func AsMessageType(b []byte) (MessageType, error) {
	c := NewMessageType()
	err := msgpack.Unmarshal(b, &c)
	if err != nil {
		return NewMessageType(), err
	}
	return c, nil
}

// MessageTypeAsBytes - упаковывает объект в массив байтов
func MessageTypeAsBytes(m *MessageType) ([]byte, error) {
	b, err := msgpack.Marshal(m)
	if err != nil {
		return nil, err
	}
	return b, nil
}

// GetStructVersion - возвращает версию модели
func (m MessageType) GetStructVersion() uint32 {
	if versionMessageType == 0 {
		versionMessageType = CalcStructVersion(reflect.TypeOf(m))
	}

	return versionMessageType
}

// GetModelFromJSON - создаёт модель из строки json
func (m *MessageType) GetModelFromJSON(sModel string) error {
	var err error

	var bytes []byte
	bytes = []byte(sModel)

	err = json.Unmarshal(bytes, m)

	return err
}

// GetJSON - возвращает строку json из модели
func (m MessageType) GetJSON() (string, error) {
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
func (m *MessageType) Read() error {
	err := crud_MessageType.Read(m)

	return err
}

// Save - записывает объект в БД по ID
func (m *MessageType) Save() error {
	err := crud_MessageType.Save(m)

	return err
}

// Update - обновляет объект в БД по ID
func (m *MessageType) Update() error {
	err := crud_MessageType.Update(m)

	return err
}

// Create - создаёт объект в БД с новым ID
func (m *MessageType) Create() error {
	err := crud_MessageType.Create(m)

	return err
}

// Delete - устанавливает признак пометки удаления в БД
func (m *MessageType) Delete() error {
	err := crud_MessageType.Delete(m)

	return err
}

// Restore - снимает признак пометки удаления в БД
func (m *MessageType) Restore() error {
	err := crud_MessageType.Restore(m)

	return err
}

// SetCrudInterface - заполняет интерфейс crud: DB, GRPC, NRPC
func (m MessageType) SetCrudInterface(crud ICrud_MessageType) {
	crud_MessageType = crud

	return
}

//---------------------------- конец CRUD операции ------------------------------------------------------------
