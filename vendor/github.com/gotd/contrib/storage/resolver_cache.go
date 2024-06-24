package storage

import (
	"context"

	"github.com/go-faster/errors"

	"github.com/gotd/td/telegram/message/peer"
	"github.com/gotd/td/tg"
)

// ResolverCache is a peer.Resolver cache implemented using peer storage.
type ResolverCache struct {
	next    peer.Resolver
	storage PeerStorage
}

// NewResolverCache creates new ResolverCache.
func NewResolverCache(next peer.Resolver, storage PeerStorage) ResolverCache {
	return ResolverCache{next: next, storage: storage}
}

func (r ResolverCache) notFound(
	ctx context.Context,
	key string,
	f func(context.Context, string) (tg.InputPeerClass, error),
) (_ tg.InputPeerClass, rerr error) {
	// If key not found, try to resolve.
	resolved, err := f(ctx, key)
	if err != nil {
		return nil, err
	}

	var value Peer
	if err := value.FromInputPeer(resolved); err != nil {
		return nil, errors.Errorf("extract object: %w", err)
	}

	if err := r.storage.Assign(ctx, key, value); err != nil {
		return nil, errors.Errorf("assign %q: %w", key, err)
	}

	return resolved, nil
}

func (r ResolverCache) tryResolve(
	ctx context.Context,
	key string,
	f func(context.Context, string) (tg.InputPeerClass, error),
) (tg.InputPeerClass, error) {
	b, err := r.storage.Resolve(ctx, key)
	if err != nil {
		if errors.Is(err, ErrPeerNotFound) {
			return r.notFound(ctx, key, f)
		}
		return nil, errors.Errorf("get %q: %w", key, err)
	}
	return b.AsInputPeer(), nil
}

// ResolveDomain implements peer.Resolver
func (r ResolverCache) ResolveDomain(ctx context.Context, domain string) (tg.InputPeerClass, error) {
	return r.tryResolve(ctx, domain, r.next.ResolveDomain)
}

// ResolvePhone implements peer.Resolver
func (r ResolverCache) ResolvePhone(ctx context.Context, phone string) (tg.InputPeerClass, error) {
	return r.tryResolve(ctx, phone, r.next.ResolvePhone)
}
