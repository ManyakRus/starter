package object_model

import (
	"encoding/json"
	"reflect"
	"time"
)

// versionOrganizationCasebook - версия структуры модели, с учётом имен и типов полей
var versionOrganizationCasebook uint32

// crud_OrganizationCasebook - объект контроллер crud операций
var crud_OrganizationCasebook ICrud_OrganizationCasebook

type OrganizationCasebook struct {
	CommonStruct
	INN        string `json:"inn"             gorm:"column:inn;default:\"\""`
	JSONFileID int64  `json:"json_file_id"    gorm:"column:json_file_id;default:null"`
	//KPP            string    `json:"kpp"             gorm:"column:kpp;default:\"\""`
	OrganizationID int64     `json:"organization_id" gorm:"column:organization_id;default:null"`
	PDFFileID      int64     `json:"pdf_file_id"     gorm:"column:pdf_file_id;default:null"`
	JSONUpdatedAt  time.Time `json:"json_updated_at"      gorm:"column:json_updated_at;default:null"`
	PDFUpdatedAt   time.Time `json:"pdf_updated_at"      gorm:"column:pdf_updated_at;default:null"`
}

type ICrud_OrganizationCasebook interface {
	Read(o *OrganizationCasebook) error
	Save(o *OrganizationCasebook) error
	Update(o *OrganizationCasebook) error
	Create(o *OrganizationCasebook) error
	Delete(o *OrganizationCasebook) error
	Restore(o *OrganizationCasebook) error
	//Find_ByInnKpp(o *OrganizationCasebook) error
	Find_ByInn(o *OrganizationCasebook) error
	Find_ByOrganizationId(o *OrganizationCasebook) error
}

// TableName - возвращает имя таблицы в БД, нужен для gorm
func (o OrganizationCasebook) TableNameDB() string {
	return "organization_casebook"
}

// GetID - возвращает ID объекта
func (o OrganizationCasebook) GetID() int64 {
	return o.ID
}

// GetStructVersion - возвращает версию модели
func (o OrganizationCasebook) GetStructVersion() uint32 {
	if versionOrganizationCasebook == 0 {
		versionOrganizationCasebook = CalcStructVersion(reflect.TypeOf(o))
	}

	return versionOrganizationCasebook
}

// GetModelFromJSON - заполняет Модель из строки JSON
func (o *OrganizationCasebook) GetModelFromJSON(sModel string) error {
	var err error

	var bytes []byte
	bytes = []byte(sModel)

	err = json.Unmarshal(bytes, o)

	return err
}

// GetJSON - возвращает строку json из модели
func (o OrganizationCasebook) GetJSON() (string, error) {
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
func (o *OrganizationCasebook) Read() error {
	err := crud_OrganizationCasebook.Read(o)

	return err
}

// Save - записывает объект в БД по ID
func (o *OrganizationCasebook) Save() error {
	err := crud_OrganizationCasebook.Save(o)

	return err
}

// Update - обновляет объект в БД по ID
func (o *OrganizationCasebook) Update() error {
	err := crud_OrganizationCasebook.Update(o)

	return err
}

// Create - создаёт объект в БД с новым ID
func (o *OrganizationCasebook) Create() error {
	err := crud_OrganizationCasebook.Create(o)

	return err
}

// Delete - устанавливает признак пометки удаления в БД
func (o *OrganizationCasebook) Delete() error {
	err := crud_OrganizationCasebook.Delete(o)

	return err
}

// Restore - снимает признак пометки удаления в БД
func (o *OrganizationCasebook) Restore() error {
	err := crud_OrganizationCasebook.Restore(o)

	return err
}

//// Find_ByInnKpp - находит запись по ИНН и КПП
//func (o *OrganizationCasebook) Find_ByInnKpp() error {
//	err := crud_OrganizationCasebook.Find_ByInnKpp(o)
//
//	return err
//}

// Find_ByInnKpp - находит запись по ИНН и КПП
func (o *OrganizationCasebook) Find_ByInn() error {
	err := crud_OrganizationCasebook.Find_ByInn(o)

	return err
}

// Find_ByOrganizationId - находит запись по OrganizationId
func (o *OrganizationCasebook) Find_ByOrganizationId() error {
	err := crud_OrganizationCasebook.Find_ByOrganizationId(o)

	return err
}

// SetCrudInterface - заполняет интерфейс crud: DB, GRPC, NRPC
func (o OrganizationCasebook) SetCrudInterface(crud ICrud_OrganizationCasebook) {
	crud_OrganizationCasebook = crud

	return
}

//---------------------------- конец CRUD операции ------------------------------------------------------------
