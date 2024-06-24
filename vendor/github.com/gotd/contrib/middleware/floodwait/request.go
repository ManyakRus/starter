package floodwait

import (
	"context"

	"github.com/gotd/td/bin"
	"github.com/gotd/td/tg"
)

// object is a abstraction for Telegram API object with TypeID.
type object interface {
	TypeID() uint32
}

type key uint64

func (k *key) fromEncoder(encoder bin.Encoder) {
	obj, ok := encoder.(object)
	if !ok {
		return
	}
	*k = key(obj.TypeID())
}

type request struct {
	ctx    context.Context
	input  bin.Encoder
	output bin.Decoder
	next   tg.Invoker
	key    key

	retry  int
	result chan error
}
