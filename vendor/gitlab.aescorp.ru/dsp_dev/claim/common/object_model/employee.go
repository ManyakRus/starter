package object_model

import (
	"encoding/json"
	"github.com/vmihailenco/msgpack/v5"
	"reflect"
)

// versionStructEmployee - версия структуры модели, с учётом имен и типов полей
var versionStructEmployee uint32

// crud_Employee - объект контроллер crud операций
var crud_Employee ICrud_Employee

// Employee Сотрудники (Справочник).
type Employee struct {
	CommonStruct
	NameStruct
	GroupStruct
	BranchID     int64  `json:"branch_id"       gorm:"column:branch_id;default:null"`
	Email        string `json:"email"           gorm:"column:email;default:\"\""`
	IsActive     bool   `json:"is_active"       gorm:"column:is_active;default:false"`
	Login        string `json:"login"           gorm:"column:login;default:\"\""`
	ParentName   string `json:"parent_name"     gorm:"column:parent_name;default:\"\""`
	Phone        string `json:"phone"           gorm:"column:phone;default:\"\""`
	Photo        string `json:"photo"           gorm:"column:photo;default:\"\""`
	Position     string `json:"position"        gorm:"column:position;default:\"\""`
	SecondName   string `json:"second_name"     gorm:"column:second_name;default:\"\""`
	Tag          string `json:"tag"             gorm:"column:tag;default:\"\""`
	ConnectionID int64  `json:"connection_id"   gorm:"column:connection_id;default:null"`
}

type ICrud_Employee interface {
	Read(e *Employee) error
	Save(e *Employee) error
	Update(e *Employee) error
	Create(e *Employee) error
	Delete(e *Employee) error
	Restore(e *Employee) error
	Find_ByExtID(e *Employee) error
	Find_ByLogin(e *Employee) error
	Find_ByEMail(e *Employee) error
	Find_ByFIO(e *Employee) error
}

// TableName - возвращает имя таблицы в БД, нужен для gorm
func (c Employee) TableNameDB() string {
	return "employees"
}

// GetID - возвращает ID объекта
func (c Employee) GetID() int64 {
	return c.ID
}

// NewEmployee Сотрудник
func NewEmployee() Employee {
	return Employee{}
}

func AsEmployee(b []byte) (Employee, error) {
	e := NewEmployee()
	err := msgpack.Unmarshal(b, &e)
	if err != nil {
		return NewEmployee(), err
	}
	return e, nil
}

func EmployeeAsBytes(e *Employee) ([]byte, error) {
	b, err := msgpack.Marshal(e)
	if err != nil {
		return nil, err
	}
	return b, nil
}

// GetStructVersion - возвращает версию модели
func (e Employee) GetStructVersion() uint32 {
	if versionStructEmployee == 0 {
		versionStructEmployee = CalcStructVersion(reflect.TypeOf(e))
	}

	return versionStructEmployee
}

// GetModelFromJSON - заполняет Модель из строки JSON
func (e *Employee) GetModelFromJSON(sModel string) error {
	var err error

	var bytes []byte
	bytes = []byte(sModel)

	err = json.Unmarshal(bytes, e)

	return err
}

// GetJSON - возвращает строку json из модели
func (e Employee) GetJSON() (string, error) {
	var ReturnVar string
	var err error

	bytes, err := json.Marshal(e)
	if err != nil {
		return ReturnVar, err
	}
	ReturnVar = string(bytes)
	return ReturnVar, err
}

//---------------------------- CRUD операции ------------------------------------------------------------

// Read - находит запись в БД по ID, и заполняет в объект
func (e *Employee) Read() error {
	err := crud_Employee.Read(e)

	return err
}

// Save - записывает объект в БД по ID
func (e *Employee) Save() error {
	err := crud_Employee.Save(e)

	return err
}

// Update - обновляет объект в БД по ID
func (e *Employee) Update() error {
	err := crud_Employee.Update(e)

	return err
}

// Create - создаёт объект в БД с новым ID
func (e *Employee) Create() error {
	err := crud_Employee.Create(e)

	return err
}

// Delete - устанавливает признак пометки удаления в БД
func (e *Employee) Delete() error {
	err := crud_Employee.Delete(e)

	return err
}

// Restore - снимает признак пометки удаления в БД
func (e *Employee) Restore() error {
	err := crud_Employee.Restore(e)

	return err
}

// Find_ByExtID - находит объект по ExtID
func (e *Employee) Find_ByExtID() error {
	err := crud_Employee.Find_ByExtID(e)

	return err
}

// Find_ByEMail - находит объект по email
func (e *Employee) Find_ByEMail() error {
	err := crud_Employee.Find_ByEMail(e)

	return err
}

// Find_ByLogin - находит объект по Login
func (e *Employee) Find_ByLogin() error {
	err := crud_Employee.Find_ByLogin(e)

	return err
}

// Find_ByFIO - находит объект по ФИО
func (e *Employee) Find_ByFIO() error {
	err := crud_Employee.Find_ByFIO(e)

	return err
}

// SetCrudInterface - заполняет интерфейс crud: DB, GRPC, NRPC
func (e Employee) SetCrudInterface(crud ICrud_Employee) {
	crud_Employee = crud

	return
}

func (e Employee) SetCrud_Transport(i ICrud_Employee) {
	crud_Employee = i
}

//---------------------------- конец CRUD операции ------------------------------------------------------------
