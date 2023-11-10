package object_model

import (
	"encoding/json"
	"reflect"
	"time"

	"github.com/vmihailenco/msgpack/v5"
)

// versionStructOrganization - версия структуры модели, с учётом имен и типов полей
var versionStructOrganization uint32

// crud_Organization - объект контроллер crud операций
var crud_Organization ICrud_Organization

// Organization Юридическое лицо (справочник).
type Organization struct {
	CommonStruct
	NameStruct
	GroupStruct
	BankruptAt     time.Time             `json:"bankrupt_at"     gorm:"column:bankrupt_at"`
	BookkeeperName string                `json:"bookkeeper_name" gorm:"column:bookkeeper_name;default:\"\""`
	CategoryID     int64                 `json:"category_id"     gorm:"column:category_id;default:null"`
	ConnectionID   int64                 `json:"connection_id"   gorm:"column:connection_id;default:null"`
	Email          string                `json:"email"           gorm:"column:email;default:\"\""`
	FullName       string                `json:"full_name"       gorm:"column:full_name;default:\"\""`
	INN            string                `json:"inn"             gorm:"column:inn;default:\"\""`
	IsActive       bool                  `json:"is_active"       gorm:"column:is_active;default:false"`
	IsBankrupt     bool                  `json:"is_bankrupt"     gorm:"column:is_bankrupt;default:false"`
	IsLiquidated   bool                  `json:"is_liquidated"   gorm:"column:is_liquidated;default:false"`
	KPP            string                `json:"kpp"             gorm:"column:kpp;default:\"\""`
	LegalAddress   string                `json:"legal_address"   gorm:"column:legal_address;default:\"\""`
	LegalTypeID    int64                 `json:"legal_type_id"   gorm:"column:legal_type_id;default:0"`
	LiquidateAt    time.Time             `json:"liquidate_at"    gorm:"column:liquidate_at"`
	ManagerName    string                `json:"manager_name"    gorm:"column:manager_name;default:\"\""`
	NSIFlat        string                `json:"nsi_flat"        gorm:"column:nsi_flat;default:\"\""` // Значение квартиры из НСИ
	NSIFlatID      int64                 `json:"nsi_flat_id"     gorm:"column:nsi_flat_id;default:0"` // ИД типа квартиры из НСИ
	NSIID          int64                 `json:"nsi_id"          gorm:"column:nsi_id;default:0"`      // ИД адреса из НСИ
	OGRN           string                `json:"ogrn"            gorm:"column:ogrn;default:\"\""`
	OKATO          string                `json:"okato"           gorm:"column:okato;default:\"\""`
	OKPO           string                `json:"okpo"            gorm:"column:okpo;default:\"\""`
	Phone          string                `json:"phone"           gorm:"column:phone;default:\"\""`
	PostAddress    string                `json:"post_address"    gorm:"column:post_address;default:\"\""`
	RegistrationAt time.Time             `json:"registration_at" gorm:"column:registration_at;default:null"`
	State          OrganizationStateType `json:"state"           gorm:"-:all"`                          // Статус организации из НСИ.
	StateCode      string                `json:"state_code"      gorm:"column:state_code;default:\"\""` // Код статуса организации из НСИ.
	StateID        int64                 `json:"state_id"        gorm:"column:state_id;default:null"`   // ID статуса организации из НСИ.
	WWW            string                `json:"www"             gorm:"column:www;default:\"\""`

	// LegalType      LegalType             `json:"legal_type"      gorm:"-:all"` // TODO LegalType

	Accounts []Account `json:"accounts"        gorm:"-:all"`
}

type ICrud_Organization interface {
	Read(o *Organization) error
	Save(o *Organization) error
	Update(o *Organization) error
	Create(o *Organization) error
	Delete(o *Organization) error
	Restore(o *Organization) error
	Find_ByExtID(o *Organization) error
	Find_ByInnKpp(o *Organization) error
}

// TableName - возвращает имя таблицы в БД, нужен для gorm
func (o Organization) TableNameDB() string {
	return "organizations"
}

// GetID - возвращает ID объекта
func (o Organization) GetID() int64 {
	return o.ID
}

// NewOrganization -
func NewOrganization() Organization {
	return Organization{}
}

func AsOrganization(b []byte) (Organization, error) {
	c := NewOrganization()
	err := msgpack.Unmarshal(b, &c)
	if err != nil {
		return NewOrganization(), err
	}
	return c, nil
}

func OrganizationAsBytes(o *Organization) ([]byte, error) {
	b, err := msgpack.Marshal(o)
	if err != nil {
		return nil, err
	}
	return b, nil
}

// GetStructVersion - возвращает версию модели
func (o Organization) GetStructVersion() uint32 {
	if versionStructOrganization == 0 {
		versionStructOrganization = CalcStructVersion(reflect.TypeOf(o))
	}

	return versionStructOrganization
}

// GetModelFromJSON - заполняет Модель из строки JSON
func (o *Organization) GetModelFromJSON(sModel string) error {
	var err error

	var bytes []byte
	bytes = []byte(sModel)

	err = json.Unmarshal(bytes, o)

	return err
}

// GetJSON - возвращает строку json из модели
func (o Organization) GetJSON() (string, error) {
	var ReturnVar string
	var err error

	bytes, err := json.Marshal(o)
	if err != nil {
		return ReturnVar, err
	}
	ReturnVar = string(bytes)
	return ReturnVar, err
}

//---------------------------- CRUD операции ------------------------------------------------------------

// Read - находит запись в БД по ID, и заполняет в объект
func (o *Organization) Read() error {
	err := crud_Organization.Read(o)

	return err
}

// Save - записывает объект в БД по ID
func (o *Organization) Save() error {
	err := crud_Organization.Save(o)

	return err
}

// Update - обновляет объект в БД по ID
func (o *Organization) Update() error {
	err := crud_Organization.Update(o)

	return err
}

// Create - создаёт объект в БД с новым ID
func (o *Organization) Create() error {
	err := crud_Organization.Create(o)

	return err
}

// Delete - устанавливает признак пометки удаления в БД
func (o *Organization) Delete() error {
	err := crud_Organization.Delete(o)

	return err
}

// Restore - снимает признак пометки удаления в БД
func (o *Organization) Restore() error {
	err := crud_Organization.Restore(o)

	return err
}

// Find_ByExtID - находит запись по ExtID
func (o *Organization) Find_ByExtID() error {
	err := crud_Organization.Find_ByExtID(o)

	return err
}

// Find_ByInnKpp - находит запись по ИНН и КПП
// если передаётся пустой КПП, то ищет без учёта КПП
func (o *Organization) Find_ByInnKpp() error {
	err := crud_Organization.Find_ByInnKpp(o)

	return err
}

// SetCrudInterface - заполняет интерфейс crud: DB, GRPC, NRPC
func (o Organization) SetCrudInterface(crud ICrud_Organization) {
	crud_Organization = crud

	return
}

//---------------------------- конец CRUD операции ------------------------------------------------------------
