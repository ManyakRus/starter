package storage

import (
	"context"

	"github.com/go-faster/errors"

	"github.com/gotd/td/telegram/query/channels/participants"
	"github.com/gotd/td/telegram/query/dialogs"
	"github.com/gotd/td/tg"
)

// PeerCollector is a simple helper to collect peers from different sources.
type PeerCollector struct {
	storage PeerStorage
}

// Dialogs collects peers from dialog iterator.
func (c PeerCollector) Dialogs(ctx context.Context, iter *dialogs.Iterator) error {
	for iter.Next(ctx) {
		var (
			p     Peer
			value = iter.Value()
		)
		switch dlg := value.Dialog.GetPeer().(type) {
		case *tg.PeerUser:
			user, ok := value.Entities.User(dlg.UserID)
			if !ok || !p.FromUser(user) {
				continue
			}
		case *tg.PeerChat:
			chat, ok := value.Entities.Chat(dlg.ChatID)
			if !ok || !p.FromChat(chat) {
				continue
			}
		case *tg.PeerChannel:
			channel, ok := value.Entities.Channel(dlg.ChannelID)
			if !ok || !p.FromChat(channel) {
				continue
			}
		}

		if err := c.storage.Add(ctx, p); err != nil {
			return errors.Errorf("add: %w", err)
		}
	}

	return iter.Err()
}

// Participants collects peers from participants iterator.
func (c PeerCollector) Participants(ctx context.Context, iter *participants.Iterator) error {
	for iter.Next(ctx) {
		var (
			p     Peer
			value = iter.Value()
		)
		user, ok := value.User()
		if !ok {
			continue
		}

		if !p.FromUser(user) {
			continue
		}
		if err := c.storage.Add(ctx, p); err != nil {
			return errors.Errorf("add: %w", err)
		}
	}

	return iter.Err()
}

// Contacts collects peers from contacts iterator.
func (c PeerCollector) Contacts(ctx context.Context, contacts *tg.ContactsContacts) error {
	for _, user := range contacts.Users {
		var p Peer
		if !p.FromUser(user) {
			continue
		}
		if err := c.storage.Add(ctx, p); err != nil {
			return errors.Errorf("add: %w", err)
		}
	}
	return nil
}

// CollectPeers creates new PeerCollector.
func CollectPeers(storage PeerStorage) PeerCollector {
	return PeerCollector{storage: storage}
}
