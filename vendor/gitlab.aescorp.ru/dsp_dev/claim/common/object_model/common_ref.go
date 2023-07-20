package object_model

import (
	"github.com/vmihailenco/msgpack/v5"
)

// CommonRef Справочники
type CommonRef struct {
	Banks            []Bank              `json:"banks"`             // 1. Банки
	Branches         []Branch            `json:"branches"`          // 2. Отделения
	Courts           []Court             `json:"courts"`            // 3. Суды
	FileTemplates    []FileTemplate      `json:"file_templates"`    // 4. Шаблоны документов
	ServiceProviders []ServiceProvider   `json:"service_providers"` // 5. Поставщики услуг
	UserRoles        []UserRole          `json:"user_roles"`        // 6. Роли сотрудников
	WhiteList        []ContractWhiteItem `json:"white_list"`        // 7. Белый список
	BlackList        []ContractBlackItem `json:"black_list"`        // 8. Чёрный список
}

// NewCommonRef Новый набор глобальных справочников
func NewCommonRef() CommonRef {
	return CommonRef{}
}
func AsCommonRef(b []byte) (CommonRef, error) {
	c := NewCommonRef()
	err := msgpack.Unmarshal(b, &c)
	if err != nil {
		return NewCommonRef(), err
	}
	return c, nil
}

func CommonRefAsBytes(c *CommonRef) ([]byte, error) {
	b, err := msgpack.Marshal(c)
	if err != nil {
		return nil, err
	}
	return b, nil
}
