package object_model

import (
	"encoding/json"
	"github.com/vmihailenco/msgpack/v5"
	"reflect"
)

// versionStructFile - версия структуры модели, с учётом имен и типов полей
var versionStructFile uint32

// crud_File - объект контроллер crud операций
var crud_File ICrud_File

// File Файлы.
type File struct {
	CommonStruct
	GroupStruct
	NameStruct
	ExtLinkStruct
	BranchID   int64  `json:"branch_id"       gorm:"column:branch_id;default:null"`
	EmployeeID int64  `json:"employee_id"     gorm:"column:employee_id;default:null"`
	Extension  string `json:"extension"       gorm:"column:extension;default:\"\""`
	FileID     string `json:"file_id"         gorm:"column:file_id;default:\"\""`
	FileName   string `json:"file_name"       gorm:"column:file_name;default:\"\""`
	FileTypeID int64  `json:"file_type_id"    gorm:"column:file_type_id;default:null"`
	FullName   string `json:"full_name"       gorm:"column:full_name;default:\"\""`
	Size       int64  `json:"size"            gorm:"column:size;default:null"`
	TemplateID int64  `json:"template_id"     gorm:"column:template_id;default:null"`
	Version    int    `json:"version"         gorm:"column:version;default:0"`
}

type ICrud_File interface {
	Read(f *File) error
	Save(f *File) error
	Update(f *File) error
	Create(f *File) error
	Delete(f *File) error
	Restore(f *File) error
	Find_ByFileId(f *File) error
	Find_ByFullName(f *File) error
	//Find_ByExtID(ext_id int64, connection_id int64) (File, error)
}

// TableName - возвращает имя таблицы в БД, нужен для gorm
func (c File) TableNameDB() string {
	return "files"
}

// GetID - возвращает ID объекта
func (c File) GetID() int64 {
	return c.ID
}

// NewFile Файл, который физически хранится в файловом хранилище
func NewFile() File {
	return File{}
}

func AsFile(b []byte) (File, error) {
	f := NewFile()
	err := msgpack.Unmarshal(b, &f)
	if err != nil {
		return NewFile(), err
	}
	return f, nil
}

func FileAsBytes(f *File) ([]byte, error) {
	b, err := msgpack.Marshal(f)
	if err != nil {
		return nil, err
	}
	return b, nil
}

// GetStructVersion - возвращает версию модели
func (f File) GetStructVersion() uint32 {
	if versionStructFile == 0 {
		versionStructFile = CalcStructVersion(reflect.TypeOf(f))
	}

	return versionStructFile
}

// GetModelFromJSON - заполняет Модель из строки JSON
func (f *File) GetModelFromJSON(sModel string) error {
	var err error

	var bytes []byte
	bytes = []byte(sModel)

	err = json.Unmarshal(bytes, f)

	return err
}

// GetJSON - возвращает строку json из модели
func (f File) GetJSON() (string, error) {
	var ReturnVar string
	var err error

	bytes, err := json.Marshal(f)
	if err != nil {
		return ReturnVar, err
	}
	ReturnVar = string(bytes)
	return ReturnVar, err
}

//---------------------------- CRUD операции ------------------------------------------------------------

// Read - находит запись в БД по ID, и заполняет в объект
func (f *File) Read() error {
	err := crud_File.Read(f)

	return err
}

// Save - записывает объект в БД по ID
func (f *File) Save() error {
	err := crud_File.Save(f)

	return err
}

// Update - обновляет объект в БД по ID
func (f *File) Update() error {
	err := crud_File.Update(f)

	return err
}

// Create - создаёт объект в БД с новым ID
func (f *File) Create() error {
	err := crud_File.Create(f)

	return err
}

// Delete - устанавливает признак пометки удаления в БД
func (f *File) Delete() error {
	err := crud_File.Delete(f)

	return err
}

// Restore - снимает признак пометки удаления в БД
func (f *File) Restore() error {
	err := crud_File.Restore(f)

	return err
}

// Find_ByFileId - находит запись по FileID
func (f *File) Find_ByFileId() error {
	err := crud_File.Find_ByFileId(f)

	return err
}

// Find_ByFull_name - находит запись по FullName
func (f *File) Find_ByFull_name() error {
	err := crud_File.Find_ByFullName(f)

	return err
}

// SetCrudInterface - заполняет интерфейс crud: DB, GRPC, NRPC
func (f File) SetCrudInterface(crud ICrud_File) {
	crud_File = crud

	return
}

//---------------------------- конец CRUD операции ------------------------------------------------------------
