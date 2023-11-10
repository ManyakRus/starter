package object_model

import (
	"encoding/json"
	"reflect"
)

// versionStructConnection - версия структуры модели, с учётом имен и типов полей
var versionStructConnection uint32

// crud_Connection - объект контроллер crud операций
var crud_Connection ICrud_Connection

type Connection struct {
	ID       int64  `json:"id"        gorm:"column:id;primaryKey;autoIncrement:true"`
	Name     string `json:"name"      gorm:"column:name;default:\"\""`
	IsLegal  bool   `json:"is_legal"  gorm:"column:is_legal;default:false"`
	BranchId int64  `json:"branch_id" gorm:"column:branch_id;default:0"`
	Server   string `json:"server"    gorm:"column:server;default:\"\""`
	Port     string `json:"port"      gorm:"column:port;default:\"\""`
	DbName   string `json:"db_name"   gorm:"column:db_name;default:\"\""`
	DbScheme string `json:"db_scheme" gorm:"column:db_scheme;default:\"\""`
	Login    string `json:"login"     gorm:"column:login;default:\"\""`
	Password string `json:"password"  gorm:"column:password;default:\"\""`
}

type ICrud_Connection interface {
	Read(c *Connection) error
	Save(c *Connection) error
	Update(c *Connection) error
	Create(c *Connection) error
}

// TableName - возвращает имя таблицы в БД, нужен для gorm
func (c Connection) TableNameDB() string {
	return "connections"
}

// GetID - возвращает ID объекта
func (c Connection) GetID() int64 {
	return c.ID
}

// GetStructVersion - возвращает версию модели
func (c Connection) GetStructVersion() uint32 {
	if versionStructConnection == 0 {
		versionStructConnection = CalcStructVersion(reflect.TypeOf(c))
	}

	return versionStructConnection
}

// GetModelFromJSON - заполняет Модель из строки JSON
func (c *Connection) GetModelFromJSON(sModel string) error {
	var err error

	var bytes []byte
	bytes = []byte(sModel)

	err = json.Unmarshal(bytes, c)

	return err
}

// GetJSON - возвращает строку json из модели
func (c Connection) GetJSON() (string, error) {
	var ReturnVar string
	var err error

	bytes, err := json.Marshal(c)
	if err != nil {
		return ReturnVar, err
	}
	ReturnVar = string(bytes)
	return ReturnVar, err
}

//---------------------------- CRUD операции ------------------------------------------------------------

// Read - находит запись в БД по ID, и заполняет в объект
func (c *Connection) Read() error {
	err := crud_Connection.Read(c)

	return err
}

// Save - записывает объект в БД по ID
func (c *Connection) Save() error {
	err := crud_Connection.Save(c)

	return err
}

// Update - обновляет объект в БД по ID
func (c *Connection) Update() error {
	err := crud_Connection.Update(c)

	return err
}

// Create - создаёт объект в БД с новым ID
func (c *Connection) Create() error {
	err := crud_Connection.Create(c)

	return err
}

// SetCrudInterface - заполняет интерфейс crud: DB, GRPC, NRPC
func (c Connection) SetCrudInterface(crud ICrud_Connection) {
	crud_Connection = crud

	return
}

func (c Connection) SetCrud_Transport(i ICrud_Connection) {
	crud_Connection = i
}

//---------------------------- конец CRUD операции ------------------------------------------------------------
