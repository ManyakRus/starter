package object_model

import (
	"encoding/json"
	"github.com/google/uuid"
	"github.com/vmihailenco/msgpack/v5"
	"reflect"
)

// versionMessageSendStatus - версия структуры модели, с учётом имен и типов полей
var versionMessageSendStatus uint32

// crud_MessageSendStatus - объект контроллер crud операций
var crud_MessageSendStatus ICrud_MessageSendStatus

// MessageSendStatus - Статусы отправки сообщений
type MessageSendStatus struct {
	CommonStruct
	NameStruct
	Code        int       `json:"code"        gorm:"column:code;default:0"`
	FormalName  string    `json:"formal_name"            gorm:"column:formal_name;default:\"\""`           //как в Notifier
	NotifierID  uuid.UUID `json:"notifier_id"            gorm:"type:uuid;column:notifier_id;default:\"\""` //ИД как в Notifier
	IsDelivered bool      `json:"is_delivered"    gorm:"column:is_delivered;default:false"`
}

type ICrud_MessageSendStatus interface {
	Read(m *MessageSendStatus) error
	Save(m *MessageSendStatus) error
	Update(m *MessageSendStatus) error
	Create(m *MessageSendStatus) error
	Delete(m *MessageSendStatus) error
	Restore(m *MessageSendStatus) error
}

// TableName - возвращает имя таблицы в БД, нужен для gorm
func (m MessageSendStatus) TableNameDB() string {
	return "message_send_statuses"
}

// GetID - возвращает ID объекта
func (m MessageSendStatus) GetID() int64 {
	return m.ID
}

// NewMessageSendStatuses - возвращает новый	объект
func NewMessageSendStatuses() MessageSendStatus {
	return MessageSendStatus{}
}

// AsMessageSendStatuses - создаёт объект из упакованного объекта в массиве байтов
func AsMessageSendStatuses(b []byte) (MessageSendStatus, error) {
	c := NewMessageSendStatuses()
	err := msgpack.Unmarshal(b, &c)
	if err != nil {
		return NewMessageSendStatuses(), err
	}
	return c, nil
}

// MessageSendStatusesAsBytes - упаковывает объект в массив байтов
func MessageSendStatusesAsBytes(m *MessageSendStatus) ([]byte, error) {
	b, err := msgpack.Marshal(m)
	if err != nil {
		return nil, err
	}
	return b, nil
}

// GetStructVersion - возвращает версию модели
func (m MessageSendStatus) GetStructVersion() uint32 {
	if versionMessageSendStatus == 0 {
		versionMessageSendStatus = CalcStructVersion(reflect.TypeOf(m))
	}

	return versionMessageSendStatus
}

// GetModelFromJSON - создаёт модель из строки json
func (m *MessageSendStatus) GetModelFromJSON(sModel string) error {
	var err error

	var bytes []byte
	bytes = []byte(sModel)

	err = json.Unmarshal(bytes, m)

	return err
}

// GetJSON - возвращает строку json из модели
func (m MessageSendStatus) GetJSON() (string, error) {
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
func (m *MessageSendStatus) Read() error {
	err := crud_MessageSendStatus.Read(m)

	return err
}

// Save - записывает объект в БД по ID
func (m *MessageSendStatus) Save() error {
	err := crud_MessageSendStatus.Save(m)

	return err
}

// Update - обновляет объект в БД по ID
func (m *MessageSendStatus) Update() error {
	err := crud_MessageSendStatus.Update(m)

	return err
}

// Create - создаёт объект в БД с новым ID
func (m *MessageSendStatus) Create() error {
	err := crud_MessageSendStatus.Create(m)

	return err
}

// Delete - устанавливает признак пометки удаления в БД
func (m *MessageSendStatus) Delete() error {
	err := crud_MessageSendStatus.Delete(m)

	return err
}

// Restore - снимает признак пометки удаления в БД
func (m *MessageSendStatus) Restore() error {
	err := crud_MessageSendStatus.Restore(m)

	return err
}

// SetCrudInterface - заполняет интерфейс crud: DB, GRPC, NRPC
func (m MessageSendStatus) SetCrudInterface(crud ICrud_MessageSendStatus) {
	crud_MessageSendStatus = crud

	return
}

//---------------------------- конец CRUD операции ------------------------------------------------------------
