package object_model

import (
	"fmt"

	"github.com/vmihailenco/msgpack/v5"
	"gitlab.aescorp.ru/dsp_dev/claim/common/object_model/object_view"
)

// ClaimWork ПИР
type ClaimWork struct {
	BriefCases []BriefCase `json:"brief_cases"`
}

// NewClaimWork Новый объект верхнего уровня, по сути содержит список портфелей.
func NewClaimWork() ClaimWork {
	return ClaimWork{}
}

func AsClaimWork(b []byte) (ClaimWork, error) {
	c := NewClaimWork()
	err := msgpack.Unmarshal(b, &c)
	if err != nil {
		return NewClaimWork(), err
	}
	return c, nil
}

func ClaimWorkAsBytes(c *ClaimWork) ([]byte, error) {
	b, err := msgpack.Marshal(c)
	if err != nil {
		return nil, err
	}
	return b, nil
}

// ClaimWorkView выборка
func ClaimWorkView(cw *ClaimWork, c *CommonRef, t *TypeRef, useFormat bool) (object_view.ViewMap, error) {
	result := make(object_view.ViewMap, 0)

	if c == nil {
		return result, nil
	}
	if len(cw.BriefCases) == 0 {
		return result, nil
	}

	for i := 0; i < len(cw.BriefCases); i++ {
		v, err := BriefCaseView(&cw.BriefCases[i], c, t, useFormat)
		if err != nil {
			return result, fmt.Errorf("ClaimWorkView, BriefCaseView[%v], Error: %v", i, err)
		}
		result.Append(fmt.Sprintf("%v", i), v)
	}

	return result, nil
}
