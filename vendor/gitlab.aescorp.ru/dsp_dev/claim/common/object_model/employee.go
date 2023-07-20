package object_model

import (
	"github.com/vmihailenco/msgpack/v5"
)

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
