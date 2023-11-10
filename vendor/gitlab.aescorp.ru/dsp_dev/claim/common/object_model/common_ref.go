package object_model

import (
	"github.com/vmihailenco/msgpack/v5"
)

// CommonRef Справочники
type CommonRef struct {
	Banks            []Bank              `json:"banks"`             // 1. Банки
	BlackList        []ContractBlackItem `json:"black_list"`        // 2. Чёрный список
	Branches         []Branch            `json:"branches"`          // 3. Отделения
	CompletedMonths  []CompletedMonth    `json:"completed_months"`  // 4. Закрытые месяцы
	Courts           []Court             `json:"courts"`            // 5. Суды
	FileTemplates    []FileTemplate      `json:"file_templates"`    // 6. Шаблоны документов
	ServiceProviders []ServiceProvider   `json:"service_providers"` // 7. Поставщики услуг
	UserRoles        []UserRole          `json:"user_roles"`        // 8. Роли сотрудников
	WhiteList        []ContractWhiteItem `json:"white_list"`        // 9. Белый список
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
