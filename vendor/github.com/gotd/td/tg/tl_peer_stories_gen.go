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

// PeerStories represents TL type `peerStories#9a35e999`.
// Stories¹ associated to a peer
//
// Links:
//  1. https://core.telegram.org/api/stories
//
// See https://core.telegram.org/constructor/peerStories for reference.
type PeerStories struct {
	// Flags, see TL conditional fields¹
	//
	// Links:
	//  1) https://core.telegram.org/mtproto/TL-combinators#conditional-fields
	Flags bin.Fields
	// The peer
	Peer PeerClass
	// If set, contains the ID of the maximum read story
	//
	// Use SetMaxReadID and GetMaxReadID helpers.
	MaxReadID int
	// Stories
	Stories []StoryItemClass
}

// PeerStoriesTypeID is TL type id of PeerStories.
const PeerStoriesTypeID = 0x9a35e999

// Ensuring interfaces in compile-time for PeerStories.
var (
	_ bin.Encoder     = &PeerStories{}
	_ bin.Decoder     = &PeerStories{}
	_ bin.BareEncoder = &PeerStories{}
	_ bin.BareDecoder = &PeerStories{}
)

func (p *PeerStories) Zero() bool {
	if p == nil {
		return true
	}
	if !(p.Flags.Zero()) {
		return false
	}
	if !(p.Peer == nil) {
		return false
	}
	if !(p.MaxReadID == 0) {
		return false
	}
	if !(p.Stories == nil) {
		return false
	}

	return true
}

// String implements fmt.Stringer.
func (p *PeerStories) String() string {
	if p == nil {
		return "PeerStories(nil)"
	}
	type Alias PeerStories
	return fmt.Sprintf("PeerStories%+v", Alias(*p))
}

// FillFrom fills PeerStories from given interface.
func (p *PeerStories) FillFrom(from interface {
	GetPeer() (value PeerClass)
	GetMaxReadID() (value int, ok bool)
	GetStories() (value []StoryItemClass)
}) {
	p.Peer = from.GetPeer()
	if val, ok := from.GetMaxReadID(); ok {
		p.MaxReadID = val
	}

	p.Stories = from.GetStories()
}

// TypeID returns type id in TL schema.
//
// See https://core.telegram.org/mtproto/TL-tl#remarks.
func (*PeerStories) TypeID() uint32 {
	return PeerStoriesTypeID
}

// TypeName returns name of type in TL schema.
func (*PeerStories) TypeName() string {
	return "peerStories"
}

// TypeInfo returns info about TL type.
func (p *PeerStories) TypeInfo() tdp.Type {
	typ := tdp.Type{
		Name: "peerStories",
		ID:   PeerStoriesTypeID,
	}
	if p == nil {
		typ.Null = true
		return typ
	}
	typ.Fields = []tdp.Field{
		{
			Name:       "Peer",
			SchemaName: "peer",
		},
		{
			Name:       "MaxReadID",
			SchemaName: "max_read_id",
			Null:       !p.Flags.Has(0),
		},
		{
			Name:       "Stories",
			SchemaName: "stories",
		},
	}
	return typ
}

// SetFlags sets flags for non-zero fields.
func (p *PeerStories) SetFlags() {
	if !(p.MaxReadID == 0) {
		p.Flags.Set(0)
	}
}

// Encode implements bin.Encoder.
func (p *PeerStories) Encode(b *bin.Buffer) error {
	if p == nil {
		return fmt.Errorf("can't encode peerStories#9a35e999 as nil")
	}
	b.PutID(PeerStoriesTypeID)
	return p.EncodeBare(b)
}

// EncodeBare implements bin.BareEncoder.
func (p *PeerStories) EncodeBare(b *bin.Buffer) error {
	if p == nil {
		return fmt.Errorf("can't encode peerStories#9a35e999 as nil")
	}
	p.SetFlags()
	if err := p.Flags.Encode(b); err != nil {
		return fmt.Errorf("unable to encode peerStories#9a35e999: field flags: %w", err)
	}
	if p.Peer == nil {
		return fmt.Errorf("unable to encode peerStories#9a35e999: field peer is nil")
	}
	if err := p.Peer.Encode(b); err != nil {
		return fmt.Errorf("unable to encode peerStories#9a35e999: field peer: %w", err)
	}
	if p.Flags.Has(0) {
		b.PutInt(p.MaxReadID)
	}
	b.PutVectorHeader(len(p.Stories))
	for idx, v := range p.Stories {
		if v == nil {
			return fmt.Errorf("unable to encode peerStories#9a35e999: field stories element with index %d is nil", idx)
		}
		if err := v.Encode(b); err != nil {
			return fmt.Errorf("unable to encode peerStories#9a35e999: field stories element with index %d: %w", idx, err)
		}
	}
	return nil
}

// Decode implements bin.Decoder.
func (p *PeerStories) Decode(b *bin.Buffer) error {
	if p == nil {
		return fmt.Errorf("can't decode peerStories#9a35e999 to nil")
	}
	if err := b.ConsumeID(PeerStoriesTypeID); err != nil {
		return fmt.Errorf("unable to decode peerStories#9a35e999: %w", err)
	}
	return p.DecodeBare(b)
}

// DecodeBare implements bin.BareDecoder.
func (p *PeerStories) DecodeBare(b *bin.Buffer) error {
	if p == nil {
		return fmt.Errorf("can't decode peerStories#9a35e999 to nil")
	}
	{
		if err := p.Flags.Decode(b); err != nil {
			return fmt.Errorf("unable to decode peerStories#9a35e999: field flags: %w", err)
		}
	}
	{
		value, err := DecodePeer(b)
		if err != nil {
			return fmt.Errorf("unable to decode peerStories#9a35e999: field peer: %w", err)
		}
		p.Peer = value
	}
	if p.Flags.Has(0) {
		value, err := b.Int()
		if err != nil {
			return fmt.Errorf("unable to decode peerStories#9a35e999: field max_read_id: %w", err)
		}
		p.MaxReadID = value
	}
	{
		headerLen, err := b.VectorHeader()
		if err != nil {
			return fmt.Errorf("unable to decode peerStories#9a35e999: field stories: %w", err)
		}

		if headerLen > 0 {
			p.Stories = make([]StoryItemClass, 0, headerLen%bin.PreallocateLimit)
		}
		for idx := 0; idx < headerLen; idx++ {
			value, err := DecodeStoryItem(b)
			if err != nil {
				return fmt.Errorf("unable to decode peerStories#9a35e999: field stories: %w", err)
			}
			p.Stories = append(p.Stories, value)
		}
	}
	return nil
}

// GetPeer returns value of Peer field.
func (p *PeerStories) GetPeer() (value PeerClass) {
	if p == nil {
		return
	}
	return p.Peer
}

// SetMaxReadID sets value of MaxReadID conditional field.
func (p *PeerStories) SetMaxReadID(value int) {
	p.Flags.Set(0)
	p.MaxReadID = value
}

// GetMaxReadID returns value of MaxReadID conditional field and
// boolean which is true if field was set.
func (p *PeerStories) GetMaxReadID() (value int, ok bool) {
	if p == nil {
		return
	}
	if !p.Flags.Has(0) {
		return value, false
	}
	return p.MaxReadID, true
}

// GetStories returns value of Stories field.
func (p *PeerStories) GetStories() (value []StoryItemClass) {
	if p == nil {
		return
	}
	return p.Stories
}

// MapStories returns field Stories wrapped in StoryItemClassArray helper.
func (p *PeerStories) MapStories() (value StoryItemClassArray) {
	return StoryItemClassArray(p.Stories)
}