package floodwait

import (
	"context"
	"time"

	"github.com/go-faster/errors"
	"go.uber.org/atomic"
	"golang.org/x/sync/errgroup"

	"github.com/gotd/td/telegram"
	"github.com/gotd/td/tg"

	"github.com/gotd/td/bin"
	"github.com/gotd/td/clock"
	"github.com/gotd/td/tgerr"
)

const (
	defaultTick       = time.Millisecond
	defaultMaxWait    = time.Minute
	defaultMaxRetries = 5
)

// Waiter is a tg.Invoker that handles flood wait errors on underlying invoker.
//
// This implementation uses a request scheduler and is more suitable for long-running
// programs with high level of concurrency and parallelism.
//
// You should use Waiter if unsure which waiter implementation to use.
//
// See SimpleWaiter for a simple timer-based implementation.
type Waiter struct {
	clock clock.Clock
	sch   *scheduler

	running    atomic.Bool
	tick       time.Duration
	maxWait    time.Duration
	maxRetries int
	onWait     func(ctx context.Context, wait FloodWait)
}

// FloodWait event.
type FloodWait struct {
	Duration time.Duration
}

// NewWaiter returns a new invoker that waits on the flood wait errors.
//
// NB: You MUST use Run method. Example:
//
//	if err := waiter.Run(ctx, func(ctx context.Context) error {
//		// Client should be started after waiter.
//		return client.Run(ctx, handler)
//	}); err != nil {
//		return errors.Wrap(err, "run client")
//	}
func NewWaiter() *Waiter {
	return &Waiter{
		clock:      clock.System,
		sch:        newScheduler(clock.System, time.Second),
		tick:       defaultTick,
		maxWait:    defaultMaxWait,
		maxRetries: defaultMaxRetries,
		onWait:     func(ctx context.Context, wait FloodWait) {},
	}
}

// WithCallback sets callback for flood wait event.
func (w *Waiter) WithCallback(f func(ctx context.Context, wait FloodWait)) *Waiter {
	w = w.clone()
	w.onWait = f
	return w
}

// clone returns a copy of the Waiter.
func (w *Waiter) clone() *Waiter {
	return &Waiter{
		clock:      w.clock,
		sch:        w.sch,
		tick:       w.tick,
		maxWait:    w.maxWait,
		maxRetries: w.maxRetries,
	}
}

// WithClock sets clock to use. Default is to use system clock.
func (w *Waiter) WithClock(c clock.Clock) *Waiter {
	w = w.clone()
	w.clock = c
	return w
}

// WithMaxWait limits wait time per attempt. Waiter will return an error if flood wait
// time exceeds that limit. Default is to wait at most a minute.
//
// To limit total wait time use a context.Context with timeout or deadline set.
func (w *Waiter) WithMaxWait(m time.Duration) *Waiter {
	w = w.clone()
	w.maxWait = m
	return w
}

// WithMaxRetries sets max number of retries before giving up. Default is to retry at most 5 times.
func (w *Waiter) WithMaxRetries(m int) *Waiter {
	w = w.clone()
	w.maxRetries = m
	return w
}

// WithTick sets gather tick interval for Waiter. Default is 1ms.
func (w *Waiter) WithTick(t time.Duration) *Waiter {
	w = w.clone()
	if t <= 0 {
		t = time.Nanosecond
	}
	w.tick = t
	return w
}

// Run runs send loop.
//
// Example:
//
//	if err := waiter.Run(ctx, func(ctx context.Context) error {
//		// Client should be started after waiter.
//		return client.Run(ctx, handler)
//	}); err != nil {
//		return errors.Wrap(err, "run client")
//	}
func (w *Waiter) Run(ctx context.Context, f func(ctx context.Context) error) (err error) {
	w.running.Store(true)
	defer w.running.Store(false)

	ctx, cancel := context.WithCancel(ctx)
	wg, ctx := errgroup.WithContext(ctx)
	wg.Go(func() error {
		defer cancel()
		return f(ctx)
	})
	wg.Go(func() error {
		ticker := w.clock.Ticker(w.tick)
		defer ticker.Stop()

		var requests []scheduled
		for {
			select {
			case <-ticker.C():
				requests = w.sch.gather(requests[:0])
				if len(requests) < 1 {
					continue
				}

				for _, s := range requests {
					ret, err := w.send(s)
					if ret {
						select {
						case s.request.result <- err:
						default:
						}
					}
				}
			case <-ctx.Done():
				return nil
			}
		}
	})

	return wg.Wait()
}

func (w *Waiter) send(s scheduled) (bool, error) {
	err := s.request.next.Invoke(s.request.ctx, s.request.input, s.request.output)

	d, ok := tgerr.AsFloodWait(err)
	if !ok {
		w.sch.nice(s.request.key)
		return true, err
	}

	// Notify about flood wait.
	w.onWait(s.request.ctx, FloodWait{
		Duration: d,
	})

	s.request.retry++

	if max := w.maxRetries; max != 0 && s.request.retry > max {
		return true, errors.Errorf("flood wait retry limit exceeded (%d > %d): %w", s.request.retry, max, err)
	}

	if d < time.Second {
		d = time.Second
	}
	if max := w.maxWait; max != 0 && d > max {
		return true, errors.Errorf("flood wait argument is too big (%v > %v): %w", d, max, err)
	}

	w.sch.flood(s.request, d)
	return false, nil
}

// Handle implements telegram.Middleware.
func (w *Waiter) Handle(next tg.Invoker) telegram.InvokeFunc {
	return func(ctx context.Context, input bin.Encoder, output bin.Decoder) error {
		if !w.running.Load() {
			// Return explicit error if waiter is not running.
			return errors.New("the Waiter middleware is not running: Run(ctx) method is not called or exited")
		}
		select {
		case err := <-w.sch.new(ctx, input, output, next):
			return err
		case <-ctx.Done():
			return ctx.Err()
		}
	}
}
