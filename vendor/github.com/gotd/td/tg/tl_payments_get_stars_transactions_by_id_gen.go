// Code generated by gotdgen, DO NOT EDIT.

package tg

import (
	"context"
	"errors"
	"fmt"
	"sort"
	"strings"

	"go.uber.org/multierr"

	"github.com/gotd/td/bin"
	"github.com/gotd/td/tdjson"
	"github.com/gotd/td/tdp"
	"github.com/gotd/td/tgerr"
)

// No-op definition for keeping imports.
var (
	_ = bin.Buffer{}
	_ = context.Background()
	_ = fmt.Stringer(nil)
	_ = strings.Builder{}
	_ = errors.Is
	_ = multierr.AppendInto
	_ = sort.Ints
	_ = tdp.Format
	_ = tgerr.Error{}
	_ = tdjson.Encoder{}
)

// PaymentsGetStarsTransactionsByIDRequest represents TL type `payments.getStarsTransactionsByID#27842d2e`.
// Obtain info about Telegram Star transactions »¹ using specific transaction IDs.
//
// Links:
//  1. https://core.telegram.org/api/stars#balance-and-transaction-history
//
// See https://core.telegram.org/method/payments.getStarsTransactionsByID for reference.
type PaymentsGetStarsTransactionsByIDRequest struct {
	// Channel or bot.
	Peer InputPeerClass
	// Transaction IDs.
	ID []InputStarsTransaction
}

// PaymentsGetStarsTransactionsByIDRequestTypeID is TL type id of PaymentsGetStarsTransactionsByIDRequest.
const PaymentsGetStarsTransactionsByIDRequestTypeID = 0x27842d2e

// Ensuring interfaces in compile-time for PaymentsGetStarsTransactionsByIDRequest.
var (
	_ bin.Encoder     = &PaymentsGetStarsTransactionsByIDRequest{}
	_ bin.Decoder     = &PaymentsGetStarsTransactionsByIDRequest{}
	_ bin.BareEncoder = &PaymentsGetStarsTransactionsByIDRequest{}
	_ bin.BareDecoder = &PaymentsGetStarsTransactionsByIDRequest{}
)

func (g *PaymentsGetStarsTransactionsByIDRequest) Zero() bool {
	if g == nil {
		return true
	}
	if !(g.Peer == nil) {
		return false
	}
	if !(g.ID == nil) {
		return false
	}

	return true
}

// String implements fmt.Stringer.
func (g *PaymentsGetStarsTransactionsByIDRequest) String() string {
	if g == nil {
		return "PaymentsGetStarsTransactionsByIDRequest(nil)"
	}
	type Alias PaymentsGetStarsTransactionsByIDRequest
	return fmt.Sprintf("PaymentsGetStarsTransactionsByIDRequest%+v", Alias(*g))
}

// FillFrom fills PaymentsGetStarsTransactionsByIDRequest from given interface.
func (g *PaymentsGetStarsTransactionsByIDRequest) FillFrom(from interface {
	GetPeer() (value InputPeerClass)
	GetID() (value []InputStarsTransaction)
}) {
	g.Peer = from.GetPeer()
	g.ID = from.GetID()
}

// TypeID returns type id in TL schema.
//
// See https://core.telegram.org/mtproto/TL-tl#remarks.
func (*PaymentsGetStarsTransactionsByIDRequest) TypeID() uint32 {
	return PaymentsGetStarsTransactionsByIDRequestTypeID
}

// TypeName returns name of type in TL schema.
func (*PaymentsGetStarsTransactionsByIDRequest) TypeName() string {
	return "payments.getStarsTransactionsByID"
}

// TypeInfo returns info about TL type.
func (g *PaymentsGetStarsTransactionsByIDRequest) TypeInfo() tdp.Type {
	typ := tdp.Type{
		Name: "payments.getStarsTransactionsByID",
		ID:   PaymentsGetStarsTransactionsByIDRequestTypeID,
	}
	if g == nil {
		typ.Null = true
		return typ
	}
	typ.Fields = []tdp.Field{
		{
			Name:       "Peer",
			SchemaName: "peer",
		},
		{
			Name:       "ID",
			SchemaName: "id",
		},
	}
	return typ
}

// Encode implements bin.Encoder.
func (g *PaymentsGetStarsTransactionsByIDRequest) Encode(b *bin.Buffer) error {
	if g == nil {
		return fmt.Errorf("can't encode payments.getStarsTransactionsByID#27842d2e as nil")
	}
	b.PutID(PaymentsGetStarsTransactionsByIDRequestTypeID)
	return g.EncodeBare(b)
}

// EncodeBare implements bin.BareEncoder.
func (g *PaymentsGetStarsTransactionsByIDRequest) EncodeBare(b *bin.Buffer) error {
	if g == nil {
		return fmt.Errorf("can't encode payments.getStarsTransactionsByID#27842d2e as nil")
	}
	if g.Peer == nil {
		return fmt.Errorf("unable to encode payments.getStarsTransactionsByID#27842d2e: field peer is nil")
	}
	if err := g.Peer.Encode(b); err != nil {
		return fmt.Errorf("unable to encode payments.getStarsTransactionsByID#27842d2e: field peer: %w", err)
	}
	b.PutVectorHeader(len(g.ID))
	for idx, v := range g.ID {
		if err := v.Encode(b); err != nil {
			return fmt.Errorf("unable to encode payments.getStarsTransactionsByID#27842d2e: field id element with index %d: %w", idx, err)
		}
	}
	return nil
}

// Decode implements bin.Decoder.
func (g *PaymentsGetStarsTransactionsByIDRequest) Decode(b *bin.Buffer) error {
	if g == nil {
		return fmt.Errorf("can't decode payments.getStarsTransactionsByID#27842d2e to nil")
	}
	if err := b.ConsumeID(PaymentsGetStarsTransactionsByIDRequestTypeID); err != nil {
		return fmt.Errorf("unable to decode payments.getStarsTransactionsByID#27842d2e: %w", err)
	}
	return g.DecodeBare(b)
}

// DecodeBare implements bin.BareDecoder.
func (g *PaymentsGetStarsTransactionsByIDRequest) DecodeBare(b *bin.Buffer) error {
	if g == nil {
		return fmt.Errorf("can't decode payments.getStarsTransactionsByID#27842d2e to nil")
	}
	{
		value, err := DecodeInputPeer(b)
		if err != nil {
			return fmt.Errorf("unable to decode payments.getStarsTransactionsByID#27842d2e: field peer: %w", err)
		}
		g.Peer = value
	}
	{
		headerLen, err := b.VectorHeader()
		if err != nil {
			return fmt.Errorf("unable to decode payments.getStarsTransactionsByID#27842d2e: field id: %w", err)
		}

		if headerLen > 0 {
			g.ID = make([]InputStarsTransaction, 0, headerLen%bin.PreallocateLimit)
		}
		for idx := 0; idx < headerLen; idx++ {
			var value InputStarsTransaction
			if err := value.Decode(b); err != nil {
				return fmt.Errorf("unable to decode payments.getStarsTransactionsByID#27842d2e: field id: %w", err)
			}
			g.ID = append(g.ID, value)
		}
	}
	return nil
}

// GetPeer returns value of Peer field.
func (g *PaymentsGetStarsTransactionsByIDRequest) GetPeer() (value InputPeerClass) {
	if g == nil {
		return
	}
	return g.Peer
}

// GetID returns value of ID field.
func (g *PaymentsGetStarsTransactionsByIDRequest) GetID() (value []InputStarsTransaction) {
	if g == nil {
		return
	}
	return g.ID
}

// PaymentsGetStarsTransactionsByID invokes method payments.getStarsTransactionsByID#27842d2e returning error if any.
// Obtain info about Telegram Star transactions »¹ using specific transaction IDs.
//
// Links:
//  1. https://core.telegram.org/api/stars#balance-and-transaction-history
//
// Possible errors:
//
//	400 PEER_ID_INVALID: The provided peer id is invalid.
//
// See https://core.telegram.org/method/payments.getStarsTransactionsByID for reference.
func (c *Client) PaymentsGetStarsTransactionsByID(ctx context.Context, request *PaymentsGetStarsTransactionsByIDRequest) (*PaymentsStarsStatus, error) {
	var result PaymentsStarsStatus

	if err := c.rpc.Invoke(ctx, request, &result); err != nil {
		return nil, err
	}
	return &result, nil
}