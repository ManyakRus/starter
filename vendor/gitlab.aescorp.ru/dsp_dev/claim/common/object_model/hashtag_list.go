package object_model

import (
	"fmt"

	"github.com/vmihailenco/msgpack/v5"
	"gitlab.aescorp.ru/dsp_dev/claim/common/object_model/object_view"
)

type HashtagList []Hashtag

// NewHashtagList Новый объект хештег
func NewHashtagList() HashtagList {
	return HashtagList{}
}

func AsHashtagList(b []byte) (HashtagList, error) {
	h := NewHashtagList()
	err := msgpack.Unmarshal(b, &h)
	if err != nil {
		return NewHashtagList(), err
	}
	return h, nil
}

func HashtagListAsBytes(h *HashtagList) ([]byte, error) {
	b, err := msgpack.Marshal(h)
	if err != nil {
		return nil, err
	}
	return b, nil
}

// HashtagListView Список хештегов
func HashtagListView(hl *HashtagList) (object_view.ViewMap, error) {
	result := object_view.ViewMap{}

	for i := 0; i < len(*hl); i++ {
		item := (*hl)[i]

		tmp := map[string]interface{}{
			"ID":   item.ID,
			"Name": item.Name,
		}
		result.Append(fmt.Sprintf("%v", i), tmp)
	}

	return result, nil
}
