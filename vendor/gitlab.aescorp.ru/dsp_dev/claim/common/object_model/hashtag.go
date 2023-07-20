package object_model

import (
	"github.com/vmihailenco/msgpack/v5"
)

type Hashtag struct {
	CommonStruct
	NameStruct
}

// NewHashtag Новый объект хештег
func NewHashtag() Hashtag {
	return Hashtag{}
}

func AsHashtag(b []byte) (Hashtag, error) {
	h := NewHashtag()
	err := msgpack.Unmarshal(b, &h)
	if err != nil {
		return NewHashtag(), err
	}
	return h, nil
}

func HashtagAsBytes(h *Hashtag) ([]byte, error) {
	b, err := msgpack.Marshal(h)
	if err != nil {
		return nil, err
	}
	return b, nil
}
