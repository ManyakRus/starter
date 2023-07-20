package object_model

import (
	"github.com/vmihailenco/msgpack/v5"
)

// TypeRef общие, как правило, фиксированные справочники
type TypeRef struct {
	ChannelTypes              []ChannelType              `json:"channel_types"`               // 1. Каналы сообщений
	ClaimTypes                []ClaimType                `json:"claim_types"`                 // 2. Типы дел
	ContractCategoryTypes     []ContractCategoryType     `json:"contract_category_types"`     // 3. Категории договоров
	DebtTypes                 []DebtType                 `json:"debt_types"`                  // 4. Типы задолженностей
	DirectionTypes            []DirectionType            `json:"direction_types"`             // 5. Направления сообщений
	DocumentLinkTypes         []DocumentLinkType         `json:"document_link_types"`         // 6. Типы связей документов
	DocumentTypes             []DocumentType             `json:"document_types"`              // 7. Типы документов
	FileTypes                 []FileType                 `json:"file_types"`                  // 8. Типы файлов
	GenderTypes               []GenderType               `json:"gender_types"`                // 9. Пол
	LawsuitReasonTypes        []LawsuitReasonType        `json:"lawsuit_reason_types"`        // 10. Причина отбора для претензии (Справочник).
	LawsuitStageTypes         []LawsuitStageType         `json:"lawsuit_stage_types"`         // 11. Этапы дел (справочник).
	LawsuitStatusTypes        []LawsuitStatusType        `json:"lawsuit_status_types"`        // 12. Статусы дел (справочник).
	LegalTypes                []LegalType                `json:"legal_types"`                 // 13. Тип юридического лица
	OrganizationCategoryTypes []OrganizationCategoryType `json:"organization_category_types"` // 14. Категории организаций
	OrganizationStateTypes    []OrganizationStateType    `json:"organization_state_types"`    // 15. Состояния организаций
	ServiceTypes              []ServiceType              `json:"service_types"`               // 16. Типы услуг
	TableNames                []TableName                `json:"table_names"`                 // 17. Имена таблиц для привязок
}

// NewTypeRef Новый набор справочников типов
func NewTypeRef() TypeRef {
	return TypeRef{}
}

func AsTypeRef(b []byte) (TypeRef, error) {
	c := NewTypeRef()
	err := msgpack.Unmarshal(b, &c)
	if err != nil {
		return NewTypeRef(), err
	}
	return c, nil
}

func TypeRefAsBytes(c *TypeRef) ([]byte, error) {
	b, err := msgpack.Marshal(c)
	if err != nil {
		return nil, err
	}
	return b, nil
}
