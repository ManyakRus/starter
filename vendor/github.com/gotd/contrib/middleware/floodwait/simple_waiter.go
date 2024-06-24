package floodwait

import (
	"context"
	"time"

	"github.com/go-faster/errors"

	"github.com/gotd/td/bin"
	"github.com/gotd/td/clock"
	"github.com/gotd/td/telegram"
	"github.com/gotd/td/tg"
	"github.com/gotd/td/tgerr"
)

// SimpleWaiter is a tg.Invoker that handles flood wait errors on underlying invoker.
//
// This implementation is more suitable for one-off tasks and programs with low level
// of concurrency and parallelism.
//
// See Waiter for a fully-blown scheduler-based flood wait handler.
type SimpleWaiter struct {
	clock clock.Clock

	maxRetries uint
	maxWait    time.Duration
}

// NewSimpleWaiter returns a new invoker that waits on the flood wait errors.
func NewSimpleWaiter() *SimpleWaiter {
	return &SimpleWaiter{
		clock: clock.System,
	}
}

// clone returns a copy of the SimpleWaiter.
func (w *SimpleWaiter) clone() *SimpleWaiter {
	return &SimpleWaiter{
		clock:      w.clock,
		maxWait:    w.maxWait,
		maxRetries: w.maxRetries,
	}
}

// WithClock sets clock to use. Default is to use system clock.
func (w *SimpleWaiter) WithClock(c clock.Clock) *SimpleWaiter {
	w = w.clone()
	w.clock = c
	return w
}

// WithMaxRetries sets max number of retries before giving up. Default is to keep retrying
// on flood wait errors indefinitely.
func (w *SimpleWaiter) WithMaxRetries(m uint) *SimpleWaiter {
	w = w.clone()
	w.maxRetries = m
	return w
}

// WithMaxWait limits wait time per attempt. SimpleWaiter will return an error if flood wait
// time exceeds that limit. Default is to wait without time limit.
//
// To limit total wait time use a context.Context with timeout or deadline set.
func (w *SimpleWaiter) WithMaxWait(m time.Duration) *SimpleWaiter {
	w = w.clone()
	w.maxWait = m
	return w
}

// Handle implements telegram.Middleware.
func (w *SimpleWaiter) Handle(next tg.Invoker) telegram.InvokeFunc {
	return func(ctx context.Context, input bin.Encoder, output bin.Decoder) error {
		var t clock.Timer

		var retries uint
		for {
			err := next.Invoke(ctx, input, output)
			if err == nil {
				return nil
			}

			d, ok := tgerr.AsFloodWait(err)
			if !ok {
				return err
			}

			retries++

			if max := w.maxRetries; max != 0 && retries > max {
				return errors.Errorf("flood wait retry limit exceeded (%d > %d): %w", retries, max, err)
			}

			if d == 0 {
				d = time.Second
			}

			if max := w.maxWait; max != 0 && d > max {
				return errors.Errorf("flood wait argument is too big (%v > %v): %w", d, max, err)
			}

			if t == nil {
				t = w.clock.Timer(d)
			} else {
				clock.StopTimer(t)
				t.Reset(d)
			}
			select {
			case <-t.C():
				continue
			case <-ctx.Done():
				clock.StopTimer(t)
				return ctx.Err()
			}
		}
	}
}
