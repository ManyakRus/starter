package object_model

import (
	"github.com/vmihailenco/msgpack/v5"
)

// Facsimile Соответствие участков ответственных и договоров
type Facsimile struct {
	CommonStruct
	Branch      string `json:"branch"      gorm:"column:branch;default:\"\""`
	Department  string `json:"department"  gorm:"column:department;default:\"\""`
	Responsible string `json:"responsible" gorm:"column:responsible;default:\"\""`
	Contract    string `json:"contract"    gorm:"column:contract;default:\"\""`
}

// NewFacsimile Данные факсимиле
func NewFacsimile() Facsimile {
	return Facsimile{}
}

func AsFacsimile(b []byte) (Facsimile, error) {
	c := NewFacsimile()
	err := msgpack.Unmarshal(b, &c)
	if err != nil {
		return NewFacsimile(), err
	}
	return c, nil
}

func FacsimileAsBytes(c *Facsimile) ([]byte, error) {
	b, err := msgpack.Marshal(c)
	if err != nil {
		return nil, err
	}
	return b, nil
}
