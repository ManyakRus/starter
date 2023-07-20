package object_model

import (
	"github.com/vmihailenco/msgpack/v5"
)

type HashtagLink struct {
	CommonStruct
	ExtLinkStruct
	HashtagID int64 `json:"hashtag_id" gorm:"column:hashtag_id;default:null"`
}

// NewHashtagLink Новый объект связи хештега
func NewHashtagLink() HashtagLink {
	return HashtagLink{}
}

func AsHashtagLink(b []byte) (HashtagLink, error) {
	hl := NewHashtagLink()
	err := msgpack.Unmarshal(b, &hl)
	if err != nil {
		return NewHashtagLink(), err
	}
	return hl, nil
}

func HashtagLinkAsBytes(hl *HashtagLink) ([]byte, error) {
	b, err := msgpack.Marshal(hl)
	if err != nil {
		return nil, err
	}
	return b, nil
}
