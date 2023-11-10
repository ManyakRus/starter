package object_model

import (
	"encoding/json"
	"github.com/vmihailenco/msgpack/v5"
	"reflect"
	"time"
)

// versionLawsuitStatusState - версия структуры модели, с учётом имен и типов полей
var versionLawsuitStatusState uint32

// crud_LawsuitStatusState - объект контроллер crud операций
var crud_LawsuitStatusState ICrud_LawsuitStatusState

// LawsuitStatusState История статусов дела.
type LawsuitStatusState struct {
	CommonStruct
	LawsuitID int64   `json:"lawsuit_id"   gorm:"column:lawsuit_id;default:null"`
	StatusID  int64   `json:"status_id"    gorm:"column:status_id;default:null"`
	TotalDebt float64 `json:"total_debt"   gorm:"column:total_debt;default:0"`
	Tag       string  `json:"tag"          gorm:"column:tag;default:\"\""`
	CommentID int64   `json:"comment_id"    gorm:"column:comment_id;default:null"`
	//	ReceivedFunds float64   `json:"received_funds"   gorm:"column:received_funds;default:0"`
	InvoiceSum  float64   `json:"invoice_sum"   gorm:"column:invoice_sum;default:0"`
	PaySum      float64   `json:"pay_sum"   gorm:"column:pay_sum;default:0"`
	MainSum     float64   `json:"main_sum"   gorm:"column:main_sum;default:0"`
	PenaltySum  float64   `json:"penalty_sum"   gorm:"column:penalty_sum;default:0"`
	PennySum    float64   `json:"penny_sum"   gorm:"column:penny_sum;default:0"`
	RestrictSum float64   `json:"restrict_sum"   gorm:"column:restrict_sum;default:0"`
	StatusAt    time.Time `json:"status_at"         gorm:"column:status_at;default:null"`
}

type ICrud_LawsuitStatusState interface {
	Read(l *LawsuitStatusState) error
	Save(l *LawsuitStatusState) error
	Update(l *LawsuitStatusState) error
	Create(l *LawsuitStatusState) error
	Delete(l *LawsuitStatusState) error
	Restore(l *LawsuitStatusState) error
	Fill_from_Lawsuit(Lawsuit_id int64, Status_id int64) error
	FindDebtSum(Lawsuit_id int64, Status_id int64) (float64, error)
}

// TableName - возвращает имя таблицы в БД, нужен для gorm
func (l LawsuitStatusState) TableNameDB() string {
	return "lawsuit_status_states"
}

// NewLawsuitStatusState - возвращает новый	объект
func NewLawsuitStatusState() LawsuitStatusState {
	return LawsuitStatusState{}
}

// AsLawsuitStatusState - создаёт объект из упакованного объекта в массиве байтов
func AsLawsuitStatusState(b []byte) (LawsuitStatusState, error) {
	c := NewLawsuitStatusState()
	err := msgpack.Unmarshal(b, &c)
	if err != nil {
		return NewLawsuitStatusState(), err
	}
	return c, nil
}

// LawsuitStatusStateAsBytes - упаковывает объект в массив байтов
func LawsuitStatusStateAsBytes(l *LawsuitStatusState) ([]byte, error) {
	b, err := msgpack.Marshal(l)
	if err != nil {
		return nil, err
	}
	return b, nil
}

// GetStructVersion - возвращает версию модели
func (l LawsuitStatusState) GetStructVersion() uint32 {
	if versionLawsuitStatusState == 0 {
		versionLawsuitStatusState = CalcStructVersion(reflect.TypeOf(l))
	}

	return versionLawsuitStatusState
}

// GetModelFromJSON - создаёт модель из строки json
func (l *LawsuitStatusState) GetModelFromJSON(sModel string) error {
	var err error

	var bytes []byte
	bytes = []byte(sModel)

	err = json.Unmarshal(bytes, l)

	return err
}

// GetJSON - возвращает строку json из модели
func (l LawsuitStatusState) GetJSON() (string, error) {
	var ReturnVar string
	var err error

	bytes, err := json.Marshal(l)
	if err != nil {
		return ReturnVar, err
	}
	ReturnVar = string(bytes)
	return ReturnVar, err
}

//---------------------------- CRUD операции ------------------------------------------------------------

// Read - находит запись в БД по ID, и заполняет в объект
func (l *LawsuitStatusState) Read() error {
	err := crud_LawsuitStatusState.Read(l)

	return err
}

// Save - записывает объект в БД по ID
func (l *LawsuitStatusState) Save() error {
	err := crud_LawsuitStatusState.Save(l)

	return err
}

// Update - обновляет объект в БД по ID
func (l *LawsuitStatusState) Update() error {
	err := crud_LawsuitStatusState.Update(l)

	return err
}

// Create - создаёт объект в БД с новым ID
func (l *LawsuitStatusState) Create() error {
	err := crud_LawsuitStatusState.Create(l)

	return err
}

// Delete - устанавливает признак пометки удаления в БД
func (l *LawsuitStatusState) Delete() error {
	err := crud_LawsuitStatusState.Delete(l)

	return err
}

// Restore - снимает признак пометки удаления в БД
func (l *LawsuitStatusState) Restore() error {
	err := crud_LawsuitStatusState.Restore(l)

	return err
}

func (l *LawsuitStatusState) Fill_from_Lawsuit(Lawsuit_id int64, Status_id int64) error {
	err := crud_LawsuitStatusState.Fill_from_Lawsuit(Lawsuit_id, Status_id)
	return err
}

func (l *LawsuitStatusState) FindDebtSum(Lawsuit_id int64, Status_id int64) (float64, error) {
	Otvet, err := crud_LawsuitStatusState.FindDebtSum(Lawsuit_id, Status_id)
	return Otvet, err
}

// SetCrudInterface - заполняет интерфейс crud: DB, GRPC, NRPC
func (l LawsuitStatusState) SetCrudInterface(crud ICrud_LawsuitStatusState) {
	crud_LawsuitStatusState = crud

	return
}

//---------------------------- конец CRUD операции ------------------------------------------------------------
