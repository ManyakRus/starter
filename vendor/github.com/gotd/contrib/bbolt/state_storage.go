package bbolt

import (
	"context"
	"encoding/binary"
	"fmt"

	bolt "go.etcd.io/bbolt"

	"github.com/gotd/td/telegram/updates"
)

func i2b(v int) []byte { b := make([]byte, 8); binary.LittleEndian.PutUint64(b, uint64(v)); return b }

func b2i(b []byte) int { return int(binary.LittleEndian.Uint64(b)) }

func i642b(v int64) []byte {
	b := make([]byte, 16)
	binary.LittleEndian.PutUint64(b, uint64(v))
	return b
}

func b2i64(b []byte) int64 { return int64(binary.LittleEndian.Uint64(b)) }

var _ updates.StateStorage = (*State)(nil)

// State is updates.StateStorage implementation using bbolt.
type State struct {
	db *bolt.DB
}

// NewStateStorage creates new state storage over bbolt.
//
// Caller is responsible for db.Close() invocation.
func NewStateStorage(db *bolt.DB) *State { return &State{db} }

func (s *State) GetState(_ context.Context, userID int64) (state updates.State, found bool, err error) {
	tx, err := s.db.Begin(false)
	if err != nil {
		return updates.State{}, false, err
	}
	defer func() { _ = tx.Rollback() }()

	user := tx.Bucket(i642b(userID))
	if user == nil {
		return updates.State{}, false, nil
	}

	stateBucket := user.Bucket([]byte("state"))
	if stateBucket == nil {
		return updates.State{}, false, nil
	}

	var (
		pts  = stateBucket.Get([]byte("pts"))
		qts  = stateBucket.Get([]byte("qts"))
		date = stateBucket.Get([]byte("date"))
		seq  = stateBucket.Get([]byte("seq"))
	)

	if pts == nil || qts == nil || date == nil || seq == nil {
		return updates.State{}, false, nil
	}

	return updates.State{
		Pts:  b2i(pts),
		Qts:  b2i(qts),
		Date: b2i(date),
		Seq:  b2i(seq),
	}, true, nil
}

func (s *State) SetState(_ context.Context, userID int64, state updates.State) error {
	return s.db.Update(func(tx *bolt.Tx) error {
		user, err := tx.CreateBucketIfNotExists(i642b(userID))
		if err != nil {
			return err
		}

		b, err := user.CreateBucketIfNotExists([]byte("state"))
		if err != nil {
			return err
		}

		check := func(e error) {
			if err != nil {
				return
			}
			err = e
		}

		check(b.Put([]byte("pts"), i2b(state.Pts)))
		check(b.Put([]byte("qts"), i2b(state.Qts)))
		check(b.Put([]byte("date"), i2b(state.Date)))
		check(b.Put([]byte("seq"), i2b(state.Seq)))
		return err
	})
}

func (s *State) SetPts(_ context.Context, userID int64, pts int) error {
	return s.db.Update(func(tx *bolt.Tx) error {
		user, err := tx.CreateBucketIfNotExists(i642b(userID))
		if err != nil {
			return err
		}

		state := user.Bucket([]byte("state"))
		if state == nil {
			return fmt.Errorf("state not found")
		}
		return state.Put([]byte("pts"), i2b(pts))
	})
}

func (s *State) SetQts(_ context.Context, userID int64, qts int) error {
	return s.db.Update(func(tx *bolt.Tx) error {
		user, err := tx.CreateBucketIfNotExists(i642b(userID))
		if err != nil {
			return err
		}

		state := user.Bucket([]byte("state"))
		if state == nil {
			return fmt.Errorf("state not found")
		}
		return state.Put([]byte("qts"), i2b(qts))
	})
}

func (s *State) SetDate(_ context.Context, userID int64, date int) error {
	return s.db.Update(func(tx *bolt.Tx) error {
		user, err := tx.CreateBucketIfNotExists(i642b(userID))
		if err != nil {
			return err
		}

		state := user.Bucket([]byte("state"))
		if state == nil {
			return fmt.Errorf("state not found")
		}
		return state.Put([]byte("date"), i2b(date))
	})
}

func (s *State) SetSeq(_ context.Context, userID int64, seq int) error {
	return s.db.Update(func(tx *bolt.Tx) error {
		user, err := tx.CreateBucketIfNotExists(i642b(userID))
		if err != nil {
			return err
		}

		state := user.Bucket([]byte("state"))
		if state == nil {
			return fmt.Errorf("state not found")
		}
		return state.Put([]byte("seq"), i2b(seq))
	})
}

func (s *State) SetDateSeq(_ context.Context, userID int64, date, seq int) error {
	return s.db.Update(func(tx *bolt.Tx) error {
		user, err := tx.CreateBucketIfNotExists(i642b(userID))
		if err != nil {
			return err
		}

		state := user.Bucket([]byte("state"))
		if state == nil {
			return fmt.Errorf("state not found")
		}
		if err := state.Put([]byte("date"), i2b(date)); err != nil {
			return err
		}
		return state.Put([]byte("seq"), i2b(seq))
	})
}

func (s *State) GetChannelPts(_ context.Context, userID, channelID int64) (pts int, found bool, err error) {
	tx, err := s.db.Begin(false)
	if err != nil {
		return 0, false, err
	}
	defer func() { _ = tx.Rollback() }()

	user := tx.Bucket(i642b(userID))
	if user == nil {
		return 0, false, nil
	}

	channels := user.Bucket([]byte("channels"))
	if channels == nil {
		return 0, false, nil
	}

	p := channels.Get(i642b(channelID))
	if p == nil {
		return 0, false, nil
	}

	return b2i(p), true, nil
}

func (s *State) SetChannelPts(_ context.Context, userID, channelID int64, pts int) error {
	return s.db.Update(func(tx *bolt.Tx) error {
		user, err := tx.CreateBucketIfNotExists(i642b(userID))
		if err != nil {
			return err
		}

		channels, err := user.CreateBucketIfNotExists([]byte("channels"))
		if err != nil {
			return err
		}
		return channels.Put(i642b(channelID), i2b(pts))
	})
}

func (s *State) ForEachChannels(ctx context.Context, userID int64, f func(ctx context.Context, channelID int64, pts int) error) error {
	return s.db.Update(func(tx *bolt.Tx) error {
		user, err := tx.CreateBucketIfNotExists(i642b(userID))
		if err != nil {
			return err
		}

		channels, err := user.CreateBucketIfNotExists([]byte("channels"))
		if err != nil {
			return err
		}

		return channels.ForEach(func(k, v []byte) error {
			return f(ctx, b2i64(k), b2i(v))
		})
	})
}
