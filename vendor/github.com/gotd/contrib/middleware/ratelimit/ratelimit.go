package ratelimit

import (
	"context"
	"errors"
	"time"

	"golang.org/x/time/rate"

	"github.com/gotd/td/bin"
	"github.com/gotd/td/clock"
	"github.com/gotd/td/telegram"
	"github.com/gotd/td/tg"
)

// RateLimiter is a tg.Invoker that throttles RPC calls on underlying invoker.
type RateLimiter struct {
	clock clock.Clock
	lim   *rate.Limiter
}

// New returns a new invoker rate limiter using lim.
func New(r rate.Limit, b int) *RateLimiter {
	return &RateLimiter{
		clock: clock.System,
		lim:   rate.NewLimiter(r, b),
	}
}

// clone returns a copy of the RateLimiter.
func (l *RateLimiter) clone() *RateLimiter {
	return &RateLimiter{
		clock: l.clock,
		lim:   l.lim,
	}
}

// WithClock sets clock to use. Default is to use system clock.
func (l *RateLimiter) WithClock(c clock.Clock) *RateLimiter {
	l = l.clone()
	l.clock = c
	return l
}

// wait blocks until rate limiter permits an event to happen. It returns an error if
// limiter’s burst size is misconfigured, the Context is canceled, or the expected
// wait time exceeds the Context’s Deadline.
func (l *RateLimiter) wait(ctx context.Context) error {
	// Check if ctx is already canceled.
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
	}

	now := l.clock.Now()

	r := l.lim.ReserveN(now, 1)
	if !r.OK() {
		// Limiter requires n <= lim.burst for each reservation.
		return errors.New("limiter's burst size must be greater than zero")
	}

	delay := r.DelayFrom(now)
	if delay == 0 {
		return nil
	}

	// Bail out earlier if we exceed context deadline. Note that
	// contexts use system time instead of mockable clock.
	deadline, ok := ctx.Deadline()
	if ok && delay > time.Until(deadline) {
		return context.DeadlineExceeded
	}

	t := l.clock.Timer(delay)
	defer clock.StopTimer(t)
	select {
	case <-t.C():
		return nil
	case <-ctx.Done():
		r.CancelAt(l.clock.Now())
		return ctx.Err()
	}
}

// Handle implements telegram.Middleware.
func (l *RateLimiter) Handle(next tg.Invoker) telegram.InvokeFunc {
	return func(ctx context.Context, input bin.Encoder, output bin.Decoder) error {
		if err := l.wait(ctx); err != nil {
			return err
		}
		return next.Invoke(ctx, input, output)
	}
}
