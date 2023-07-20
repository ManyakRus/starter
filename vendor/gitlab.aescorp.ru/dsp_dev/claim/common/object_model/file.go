package object_model

import (
	"github.com/vmihailenco/msgpack/v5"
)

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
