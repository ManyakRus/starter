package object_model

import (
	"encoding/json"
	"github.com/vmihailenco/msgpack/v5"
	"reflect"
	"time"
)

// versionMessage - версия структуры модели, с учётом имен и типов полей
var versionMessage uint32

// crud_Message - объект контроллер crud операций
var crud_Message ICrud_Message

// Message Сообщения (входящие и исходящие).
type Message struct {
	CommonStruct
	LawsuitID           int64     `json:"lawsuit_id"        gorm:"column:lawsuit_id;default:null"`                           //Дело (ИД)
	DirectionTypeID     int64     `json:"direction_type_id"        gorm:"column:direction_type_id;default:null"`             //ИД входящее или исходящее
	Result              string    `json:"result"            gorm:"column:result;default:\"\""`                               //Результат отправки сообщения (текст ошибки)
	ChannelTypeID       int64     `json:"channel_type_id"   gorm:"column:channel_type_id;default:null"`                      //Канал отправки сообщения (ИД)
	SendStatusID        int64     `json:"send_status_id"   gorm:"column:send_status_id;default:null"`                        //Статус отправки (ИД)
	Topic               string    `json:"topic"            gorm:"column:topic;default:\"\""`                                 //Тема письма
	ContactFrom         string    `json:"contact_from"            gorm:"column:contact_from;default:\"\""`                   //EMail от кого
	ContactTo           string    `json:"contact_to"            gorm:"column:contact_to;default:\"\""`                       //EMail кому
	EmployeeIDFrom      int64     `json:"employee_id_from"   gorm:"column:employee_id_from;default:null"`                    //Сотрудник от кого сообщение (ИД)
	EmployeeIDTo        int64     `json:"employee_id_to"   gorm:"column:employee_id_to;default:null"`                        //Сотрудник от кого (ИД)
	MessageFileID       int64     `json:"message_file_id"   gorm:"column:message_file_id;default:null"`                      //Файл с текстом письма (ИД)
	SentAt              time.Time `json:"sent_at"      gorm:"column:sent_at;default:null"`                                   //Время отправки сообщения
	ReceivedAt          time.Time `json:"received_at"      gorm:"column:received_at;default:null"`                           //Дата получения сообщения
	NotifierMailingCode string    `json:"notifier_mailing_code"            gorm:"column:notifier_mailing_code;default:\"\""` //mailing_id из нотификации
	MessageID           int64     `json:"message_type_id"   gorm:"column:message_type_id;default:null"`                      //Тип сообщения
	ExtCode             string    `json:"ext_code"            gorm:"column:ext_code;default:\"\""`                           //ШПИ (штрихкод)
	MessageTypeID       int64     `json:"message_type_id"   gorm:"column:message_type_id;default:null"`                      //тип сообщения

}

type ICrud_Message interface {
	Read(m *Message) error
	Save(m *Message) error
	Update(m *Message) error
	Create(m *Message) error
	Delete(m *Message) error
	Restore(m *Message) error
	FindBy_LawsuitID_MessageTypeID(m *Message) error
}

// TableName - возвращает имя таблицы в БД, нужен для gorm
func (m Message) TableNameDB() string {
	return "messages"
}

// GetID - возвращает ID объекта
func (m Message) GetID() int64 {
	return m.ID
}

// NewMessage - возвращает новый объект
func NewMessage() Message {
	return Message{}
}

// AsMessage - создаёт объект из упакованного объекта в массиве байтов
func AsMessage(b []byte) (Message, error) {
	c := NewMessage()
	err := msgpack.Unmarshal(b, &c)
	if err != nil {
		return NewMessage(), err
	}
	return c, nil
}

// MessageAsBytes - упаковывает объект в массив байтов
func MessageAsBytes(m *Message) ([]byte, error) {
	b, err := msgpack.Marshal(m)
	if err != nil {
		return nil, err
	}
	return b, nil
}

// GetStructVersion - возвращает версию модели
func (m Message) GetStructVersion() uint32 {
	if versionMessage == 0 {
		versionMessage = CalcStructVersion(reflect.TypeOf(m))
	}

	return versionMessage
}

// GetModelFromJSON - создаёт модель из строки json
func (m *Message) GetModelFromJSON(sModel string) error {
	var err error

	var bytes []byte
	bytes = []byte(sModel)

	err = json.Unmarshal(bytes, m)

	return err
}

// GetJSON - возвращает строку json из модели
func (m Message) GetJSON() (string, error) {
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
func (m *Message) Read() error {
	err := crud_Message.Read(m)

	return err
}

// Save - записывает объект в БД по ID
func (m *Message) Save() error {
	err := crud_Message.Save(m)

	return err
}

// Update - обновляет объект в БД по ID
func (m *Message) Update() error {
	err := crud_Message.Update(m)

	return err
}

// Create - создаёт объект в БД с новым ID
func (m *Message) Create() error {
	err := crud_Message.Create(m)

	return err
}

// Delete - устанавливает признак пометки удаления в БД
func (m *Message) Delete() error {
	err := crud_Message.Delete(m)

	return err
}

// Restore - снимает признак пометки удаления в БД
func (m *Message) Restore() error {
	err := crud_Message.Restore(m)

	return err
}

// FindBy_LawsuitID_MessageTypeID - находит запись в БД по lawsuit_id + message_type_id
func (m *Message) FindBy_LawsuitID_MessageTypeID() error {
	err := crud_Message.FindBy_LawsuitID_MessageTypeID(m)

	return err
}

// SetCrudInterface - заполняет интерфейс crud: DB, GRPC, NRPC
func (m Message) SetCrudInterface(crud ICrud_Message) {
	crud_Message = crud

	return
}

//---------------------------- конец CRUD операции ------------------------------------------------------------
