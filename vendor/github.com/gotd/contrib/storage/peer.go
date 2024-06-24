package storage

import (
	"encoding/json"
	"strconv"
	"strings"
	"time"

	"github.com/go-faster/errors"
	"github.com/go-faster/jx"

	"github.com/gotd/td/bin"
	"github.com/gotd/td/telegram/query/dialogs"
	"github.com/gotd/td/tg"
)

// ErrPeerUnmarshalMustInvalidate means that persisted Peer is outdated and must be invalidated.
var ErrPeerUnmarshalMustInvalidate = errors.New("outdated data for Peer (cache miss, must invalidate)")

// LatestVersion is a latest supported version of data.
const LatestVersion = 2

// Peer is abstraction for persisted peer object.
//
// Note: unmarshal error ErrPeerUnmarshalMustInvalidate MUST be considered as cache miss
// and cache entry MUST be invalidated.
//
// The only valid way to marshal and unmarshal Peer is to use UnmarshalJSON, MarshalJSON.
type Peer struct {
	Version   int
	Key       dialogs.DialogKey
	CreatedAt time.Time
	User      *tg.User
	Chat      *tg.Chat
	Channel   *tg.Channel
	Metadata  map[string]any
}

func (p Peer) String() string {
	var b strings.Builder
	switch p.Key.Kind {
	case dialogs.Chat:
		b.WriteString("Chat")
	case dialogs.Channel:
		b.WriteString("Channel")
	case dialogs.User:
		b.WriteString("User")
	}
	b.WriteString("(")
	b.WriteString(strconv.FormatInt(p.Key.ID, 10))
	b.WriteString(")")

	b.WriteString("[")
	var entities []string
	if p.User != nil {
		entities = append(entities, "User")
	}
	if p.Chat != nil {
		entities = append(entities, "Chat")
	}
	if p.Channel != nil {
		entities = append(entities, "Channel")
	}
	b.WriteString(strings.Join(entities, ", "))
	b.WriteString("]")

	return b.String()
}

func decodeObject(d *jx.Decoder, v bin.Decoder) error {
	data, err := d.Base64()
	if err != nil {
		return errors.Wrap(err, "base64")
	}
	b := &bin.Buffer{
		Buf: data,
	}
	if err := v.Decode(b); err != nil {
		return errors.Wrap(err, "decode")
	}
	return nil
}

func (p *Peer) UnmarshalJSON(data []byte) error {
	return p.Unmarshal(jx.DecodeBytes(data))
}

func (p *Peer) Unmarshal(d *jx.Decoder) error {
	var version int
	if err := d.Capture(func(d *jx.Decoder) error {
		return d.Obj(func(d *jx.Decoder, key string) error {
			if key != "Version" {
				return d.Skip()
			}
			v, err := d.Int()
			if err != nil {
				return errors.Wrap(err, "version")
			}
			version = v
			return nil
		})
	}); err != nil {
		return errors.Wrap(err, "check version")
	}
	if version != LatestVersion {
		// Ignoring.
		return ErrPeerUnmarshalMustInvalidate
	}

	// Reset.
	p.Metadata = nil
	p.User = nil
	p.Chat = nil
	p.Channel = nil
	p.CreatedAt = time.Time{}

	if err := d.Obj(func(d *jx.Decoder, key string) error {
		switch key {
		case "Version":
			v, err := d.Int()
			if err != nil {
				return errors.Wrap(err, "version")
			}
			p.Version = v
			return nil
		case "CreatedAt":
			v, err := d.Int64()
			if err != nil {
				return errors.Wrap(err, "created_at")
			}
			p.CreatedAt = time.Unix(v, 0)
			return nil
		case "Key":
			return d.Obj(func(d *jx.Decoder, key string) error {
				switch key {
				case "Kind":
					v, err := d.Int()
					if err != nil {
						return errors.Wrap(err, "kind")
					}
					p.Key.Kind = dialogs.PeerKind(v)
				case "ID":
					v, err := d.Int64()
					if err != nil {
						return errors.Wrap(err, "id")
					}
					p.Key.ID = v
				case "AccessHash":
					v, err := d.Int64()
					if err != nil {
						return errors.Wrap(err, "access_hash")
					}
					p.Key.AccessHash = v
				default:
					return d.Skip()
				}
				return nil
			})
		case "Metadata":
			var metadata map[string]any
			buf, err := d.Raw()
			if err != nil {
				return errors.Wrap(err, "metadata")
			}
			if err := json.Unmarshal(buf, &metadata); err != nil {
				return errors.Wrap(err, "unmarshal metadata")
			}
			p.Metadata = metadata
			return nil
		case "User":
			var user tg.User
			if err := decodeObject(d, &user); err != nil {
				return errors.Wrap(err, "user")
			}
			p.User = &user
			return nil
		case "Chat":
			var chat tg.Chat
			if err := decodeObject(d, &chat); err != nil {
				return errors.Wrap(err, "chat")
			}
			p.Chat = &chat
			return nil
		case "Channel":
			var channel tg.Channel
			if err := decodeObject(d, &channel); err != nil {
				return errors.Wrap(err, "channel")
			}
			p.Channel = &channel
			return nil
		default:
			return d.Skip()
		}
	}); err != nil {
		if _, ok := errors.Into[*bin.UnexpectedIDErr](err); ok {
			return ErrPeerUnmarshalMustInvalidate
		}
		return errors.Wrap(err, "decode")
	}

	return nil
}

func (p Peer) Marshal(e *jx.Encoder) error {
	type rawObject struct {
		Key   string
		Value bin.Encoder
	}
	var toEncode []rawObject
	if p.User != nil {
		toEncode = append(toEncode, rawObject{Key: "User", Value: p.User})
	}
	if p.Chat != nil {
		toEncode = append(toEncode, rawObject{Key: "Chat", Value: p.Chat})
	}
	if p.Channel != nil {
		toEncode = append(toEncode, rawObject{Key: "Channel", Value: p.Channel})
	}
	type rawValue struct {
		Key   string
		Value []byte
	}
	var values []rawValue
	for _, v := range toEncode {
		if v.Value == nil {
			continue
		}
		b := new(bin.Buffer)
		if err := v.Value.Encode(b); err != nil {
			return errors.Wrap(err, "encode")
		}
		values = append(values, rawValue{
			Key:   v.Key,
			Value: b.Buf,
		})
	}

	metadataRaw, err := json.Marshal(p.Metadata)
	if err != nil {
		return errors.Wrap(err, "marshal metadata")
	}

	e.Obj(func(e *jx.Encoder) {
		e.Field("Version", func(e *jx.Encoder) {
			e.Int(p.Version)
		})
		e.Field("Key", func(e *jx.Encoder) {
			e.Obj(func(e *jx.Encoder) {
				e.Field("Kind", func(e *jx.Encoder) {
					e.Int(int(p.Key.Kind))
				})
				e.Field("ID", func(e *jx.Encoder) {
					e.Int64(p.Key.ID)
				})
				e.Field("AccessHash", func(e *jx.Encoder) {
					e.Int64(p.Key.AccessHash)
				})
			})
		})
		e.Field("CreatedAt", func(e *jx.Encoder) {
			e.Int64(p.CreatedAt.Unix())
		})
		e.Field("Metadata", func(e *jx.Encoder) {
			e.Raw(metadataRaw)
		})
		for _, v := range values {
			e.Field(v.Key, func(e *jx.Encoder) {
				e.Base64(v.Value)
			})
		}
	})
	return nil
}

func (p Peer) MarshalJSON() ([]byte, error) {
	var e jx.Encoder
	if err := p.Marshal(&e); err != nil {
		return nil, err
	}
	return e.Bytes(), nil
}

func addIfNotEmpty(r []string, k string) []string {
	if k == "" {
		return r
	}
	return append(r, k)
}

// Keys returns list of all associated keys (phones, usernames, etc.) stored in the peer.
func (p *Peer) Keys() []string {
	// Chat does not contain usernames or phones.
	if p.Chat != nil {
		return nil
	}

	r := make([]string, 0, 4)
	switch {
	case p.User != nil:
		r = addIfNotEmpty(r, p.User.Username)
		r = addIfNotEmpty(r, p.User.Phone)
	case p.Channel != nil:
		r = addIfNotEmpty(r, p.Channel.Username)
	}

	return r
}

// FromInputPeer fills Peer object using given tg.InputPeerClass.
func (p *Peer) FromInputPeer(input tg.InputPeerClass) error {
	k := dialogs.DialogKey{}
	if err := k.FromInputPeer(input); err != nil {
		return errors.Errorf("unpack input peer: %w", err)
	}

	*p = Peer{
		Version:   LatestVersion,
		Key:       k,
		CreatedAt: time.Now(),
	}

	return nil
}

// FromChat fills Peer object using given tg.ChatClass.
func (p *Peer) FromChat(chat tg.ChatClass) bool {
	r := Peer{
		Version:   LatestVersion,
		CreatedAt: time.Now(),
	}

	switch c := chat.(type) {
	case *tg.Chat:
		r.Key.ID = c.ID
		r.Key.Kind = dialogs.Chat
		r.Chat = c
	case *tg.ChatForbidden:
		r.Key.ID = c.ID
		r.Key.Kind = dialogs.Chat
	case *tg.Channel:
		if c.Min {
			return false
		}
		r.Key.ID = c.ID
		r.Key.AccessHash = c.AccessHash
		r.Key.Kind = dialogs.Channel
		r.Channel = c
	case *tg.ChannelForbidden:
		r.Key.ID = c.ID
		r.Key.AccessHash = c.AccessHash
		r.Key.Kind = dialogs.Channel
	default:
		return false
	}

	*p = r
	return true
}

// FromUser fills Peer object using given tg.UserClass.
func (p *Peer) FromUser(user tg.UserClass) bool {
	u, ok := user.AsNotEmpty()
	if !ok {
		return false
	}

	*p = Peer{
		Version:   LatestVersion,
		CreatedAt: time.Now(),
		User:      u,
		Key: dialogs.DialogKey{
			Kind:       dialogs.User,
			ID:         u.ID,
			AccessHash: u.AccessHash,
		},
	}

	return true
}
