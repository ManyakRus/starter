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

// ConnectedBot represents TL type `connectedBot#e7e999e7`.
//
// See https://core.telegram.org/constructor/connectedBot for reference.
type ConnectedBot struct {
	// Flags field of ConnectedBot.
	Flags bin.Fields
	// CanReply field of ConnectedBot.
	CanReply bool
	// BotID field of ConnectedBot.
	BotID int64
	// Recipients field of ConnectedBot.
	Recipients BusinessRecipients
}

// ConnectedBotTypeID is TL type id of ConnectedBot.
const ConnectedBotTypeID = 0xe7e999e7

// Ensuring interfaces in compile-time for ConnectedBot.
var (
	_ bin.Encoder     = &ConnectedBot{}
	_ bin.Decoder     = &ConnectedBot{}
	_ bin.BareEncoder = &ConnectedBot{}
	_ bin.BareDecoder = &ConnectedBot{}
)

func (c *ConnectedBot) Zero() bool {
	if c == nil {
		return true
	}
	if !(c.Flags.Zero()) {
		return false
	}
	if !(c.CanReply == false) {
		return false
	}
	if !(c.BotID == 0) {
		return false
	}
	if !(c.Recipients.Zero()) {
		return false
	}

	return true
}

// String implements fmt.Stringer.
func (c *ConnectedBot) String() string {
	if c == nil {
		return "ConnectedBot(nil)"
	}
	type Alias ConnectedBot
	return fmt.Sprintf("ConnectedBot%+v", Alias(*c))
}

// FillFrom fills ConnectedBot from given interface.
func (c *ConnectedBot) FillFrom(from interface {
	GetCanReply() (value bool)
	GetBotID() (value int64)
	GetRecipients() (value BusinessRecipients)
}) {
	c.CanReply = from.GetCanReply()
	c.BotID = from.GetBotID()
	c.Recipients = from.GetRecipients()
}

// TypeID returns type id in TL schema.
//
// See https://core.telegram.org/mtproto/TL-tl#remarks.
func (*ConnectedBot) TypeID() uint32 {
	return ConnectedBotTypeID
}

// TypeName returns name of type in TL schema.
func (*ConnectedBot) TypeName() string {
	return "connectedBot"
}

// TypeInfo returns info about TL type.
func (c *ConnectedBot) TypeInfo() tdp.Type {
	typ := tdp.Type{
		Name: "connectedBot",
		ID:   ConnectedBotTypeID,
	}
	if c == nil {
		typ.Null = true
		return typ
	}
	typ.Fields = []tdp.Field{
		{
			Name:       "CanReply",
			SchemaName: "can_reply",
			Null:       !c.Flags.Has(0),
		},
		{
			Name:       "BotID",
			SchemaName: "bot_id",
		},
		{
			Name:       "Recipients",
			SchemaName: "recipients",
		},
	}
	return typ
}

// SetFlags sets flags for non-zero fields.
func (c *ConnectedBot) SetFlags() {
	if !(c.CanReply == false) {
		c.Flags.Set(0)
	}
}

// Encode implements bin.Encoder.
func (c *ConnectedBot) Encode(b *bin.Buffer) error {
	if c == nil {
		return fmt.Errorf("can't encode connectedBot#e7e999e7 as nil")
	}
	b.PutID(ConnectedBotTypeID)
	return c.EncodeBare(b)
}

// EncodeBare implements bin.BareEncoder.
func (c *ConnectedBot) EncodeBare(b *bin.Buffer) error {
	if c == nil {
		return fmt.Errorf("can't encode connectedBot#e7e999e7 as nil")
	}
	c.SetFlags()
	if err := c.Flags.Encode(b); err != nil {
		return fmt.Errorf("unable to encode connectedBot#e7e999e7: field flags: %w", err)
	}
	b.PutLong(c.BotID)
	if err := c.Recipients.Encode(b); err != nil {
		return fmt.Errorf("unable to encode connectedBot#e7e999e7: field recipients: %w", err)
	}
	return nil
}

// Decode implements bin.Decoder.
func (c *ConnectedBot) Decode(b *bin.Buffer) error {
	if c == nil {
		return fmt.Errorf("can't decode connectedBot#e7e999e7 to nil")
	}
	if err := b.ConsumeID(ConnectedBotTypeID); err != nil {
		return fmt.Errorf("unable to decode connectedBot#e7e999e7: %w", err)
	}
	return c.DecodeBare(b)
}

// DecodeBare implements bin.BareDecoder.
func (c *ConnectedBot) DecodeBare(b *bin.Buffer) error {
	if c == nil {
		return fmt.Errorf("can't decode connectedBot#e7e999e7 to nil")
	}
	{
		if err := c.Flags.Decode(b); err != nil {
			return fmt.Errorf("unable to decode connectedBot#e7e999e7: field flags: %w", err)
		}
	}
	c.CanReply = c.Flags.Has(0)
	{
		value, err := b.Long()
		if err != nil {
			return fmt.Errorf("unable to decode connectedBot#e7e999e7: field bot_id: %w", err)
		}
		c.BotID = value
	}
	{
		if err := c.Recipients.Decode(b); err != nil {
			return fmt.Errorf("unable to decode connectedBot#e7e999e7: field recipients: %w", err)
		}
	}
	return nil
}

// SetCanReply sets value of CanReply conditional field.
func (c *ConnectedBot) SetCanReply(value bool) {
	if value {
		c.Flags.Set(0)
		c.CanReply = true
	} else {
		c.Flags.Unset(0)
		c.CanReply = false
	}
}

// GetCanReply returns value of CanReply conditional field.
func (c *ConnectedBot) GetCanReply() (value bool) {
	if c == nil {
		return
	}
	return c.Flags.Has(0)
}

// GetBotID returns value of BotID field.
func (c *ConnectedBot) GetBotID() (value int64) {
	if c == nil {
		return
	}
	return c.BotID
}

// GetRecipients returns value of Recipients field.
func (c *ConnectedBot) GetRecipients() (value BusinessRecipients) {
	if c == nil {
		return
	}
	return c.Recipients
}