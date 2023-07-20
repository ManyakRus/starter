package object_model

import (
	"time"

	"github.com/vmihailenco/msgpack/v5"
	"gitlab.aescorp.ru/dsp_dev/claim/common/object_model/alias"
)

// ContractWhiteItem "Белый" список договоров. Кому не предъявляется претензия.
type ContractWhiteItem struct {
	CommonStruct
	Contract       Contract             `json:"contract"        gorm:"-:all"`
	ContractID     int64                `json:"contract_id"     gorm:"column:contract_id;default:null"`
	ContractNumber alias.ContractNumber `json:"contract_number" gorm:"column:contract_number;default:null"`
	CreatedBy      Employee             `json:"created_by"      gorm:"-:all"`
	CreatedByID    int64                `json:"created_by_id"   gorm:"column:created_by_id;default:null"`
	DateFrom       time.Time            `json:"date_from"       gorm:"column:date_from;default:null"`
	DateTo         time.Time            `json:"date_to"         gorm:"column:date_to;default:null"`
	EDMSLink       string               `json:"edms_link"       gorm:"column:edms_link;default:\"\""`
	ModifiedBy     Employee             `json:"modified_by"     gorm:"-:all"`
	ModifiedByID   int64                `json:"modified_by_id"  gorm:"column:modified_by_id;default:null"`
	Note           string               `json:"note"            gorm:"column:note;default:\"\""`
	Reason         string               `json:"reason"          gorm:"column:reason;default:\"\""`
}

// NewWhiteListItem -- Новая запись белого списка
func NewWhiteListItem() ContractWhiteItem {
	return ContractWhiteItem{}
}

func AsWhiteListItem(b []byte) (ContractWhiteItem, error) {
	c := NewWhiteListItem()
	err := msgpack.Unmarshal(b, &c)
	if err != nil {
		return NewWhiteListItem(), err
	}
	return c, nil
}

func WhiteListItemAsBytes(c *ContractWhiteItem) ([]byte, error) {
	b, err := msgpack.Marshal(c)
	if err != nil {
		return nil, err
	}
	return b, nil
}
