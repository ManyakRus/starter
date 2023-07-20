package object_model

import (
	"fmt"

	"github.com/vmihailenco/msgpack/v5"
	"gitlab.aescorp.ru/dsp_dev/claim/common/object_model/object_view"

	"gitlab.aescorp.ru/dsp_dev/claim/common/object_model/front/front_format"
)

type ContractBlackList []ContractBlackItem

// NewContractBlackList -- Новый объект чёрного списка
func NewContractBlackList() ContractBlackList {
	return ContractBlackList{}
}

func AsContractBlackList(b []byte) (ContractBlackList, error) {
	c := NewContractBlackList()
	err := msgpack.Unmarshal(b, &c)
	if err != nil {
		return NewContractBlackList(), err
	}
	return c, nil
}

func ContractBlackListAsBytes(c *ContractBlackList) ([]byte, error) {
	b, err := msgpack.Marshal(c)
	if err != nil {
		return nil, err
	}
	return b, nil
}

// BlackListView Белый список договоров
func BlackListView(cbl *ContractBlackList) (object_view.ViewMap, error) {
	result := object_view.ViewMap{}

	for i := 0; i < len(*cbl); i++ {
		item := (*cbl)[i]

		dateFrom := front_format.FrontDate(item.DateFrom)
		dateTo := front_format.FrontDate(item.DateTo)

		if item.DateFrom.Year() <= 2000 {
			dateFrom = ""
		}
		if item.DateTo.Year() >= 2100 {
			dateTo = ""
		}

		tmp := map[string]interface{}{
			"ID":                   item.ID,
			"ContractNumber":       item.Contract.Number,
			"CreatedAt":            front_format.FrontDate(item.CreatedAt),
			"CreatedBy":            item.CreatedBy.Name,
			"DateFrom":             dateFrom,
			"DateFromDatePicker":   item.DateFrom.Format("2006-01-02"),
			"DateTo":               dateTo,
			"DateToDatePicker":     item.DateTo.Format("2006-01-02"),
			"EDMSLink":             item.EDMSLink,
			"ModifiedAt":           front_format.FrontDate(item.ModifiedAt),
			"ModifiedBy":           item.ModifiedBy.Name,
			"OrganizationCategory": item.Contract.Category.Name,
			"OrganizationINN":      item.Contract.Organization.INN,
			"OrganizationKPP":      item.Contract.Organization.KPP,
			"OrganizationName":     item.Contract.Organization.Name,
			"Reason":               item.Reason,
			"IsDeleted":            item.IsDeleted,
			"DeletedAt":            front_format.FrontDate(item.DeletedAt),
		}
		result.Append(fmt.Sprintf("%v", i), tmp)
	}

	return result, nil
}
