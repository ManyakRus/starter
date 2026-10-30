// Package bg implements wrapper for running client in background.
package bg

import (
	"context"
	"errors"
	"time"
)

// Client abstracts telegram client.
type Client interface {
	Run(ctx context.Context, f func(ctx context.Context) error) error
}

// StopFunc closes Client and waits until Run returns.
type StopFunc func() error

// ErrStartupTimeout is returned by Connect when the client did not become ready
// within the configured startup timeout. See WithStartupTimeout.
var ErrStartupTimeout = errors.New("bg: client startup timeout")

type connectOptions struct {
	ctx     context.Context
	timeout time.Duration
}

// Option for Connect.
type Option interface {
	apply(o *connectOptions)
}

type fnOption func(o *connectOptions)

func (f fnOption) apply(o *connectOptions) {
	f(o)
}

// WithContext sets base context for client.
func WithContext(ctx context.Context) Option {
	return fnOption(func(o *connectOptions) {
		o.ctx = ctx
	})
}

// WithStartupTimeout bounds how long Connect waits for the client to become
// ready before giving up with ErrStartupTimeout.
//
// The client retries connection attempts indefinitely, so without a deadline
// (either via this option or a cancelable context, see WithContext) Connect can
// block forever if the connection never becomes ready. A non-positive duration
// disables the timeout, which is the default.
func WithStartupTimeout(d time.Duration) Option {
	return fnOption(func(o *connectOptions) {
		o.timeout = d
	})
}

// noopStop is returned alongside an error so callers can always invoke the
// returned StopFunc unconditionally (e.g. via defer) without a nil check.
func noopStop() error { return nil }

// Connect blocks until the client is connected and ready for requests, calling
// Run internally in background.
//
// The returned StopFunc terminates the client and waits until Run returns.
func Connect(client Client, options ...Option) (StopFunc, error) {
	opt := &connectOptions{
		ctx: context.Background(),
	}
	for _, o := range options {
		o.apply(opt)
	}

	ctx, cancel := context.WithCancel(opt.ctx)

	errC := make(chan error, 1)
	initDone := make(chan struct{})
	go func() {
		defer close(errC)
		errC <- client.Run(ctx, func(ctx context.Context) error {
			// Run invokes this callback only once the primary connection is
			// initialized and ready to send requests, so it is safe to use the
			// client right after Connect returns.
			//
			// See https://github.com/gotd/td/issues/731.
			close(initDone)
			<-ctx.Done()
			if errors.Is(ctx.Err(), context.Canceled) {
				return nil
			}
			return ctx.Err()
		})
	}()

	// Optionally bound the time we wait for readiness, otherwise Connect could
	// block forever because Run retries connection attempts indefinitely.
	var timeoutC <-chan time.Time
	if opt.timeout > 0 {
		t := time.NewTimer(opt.timeout)
		defer t.Stop()
		timeoutC = t.C
	}

	select {
	case <-ctx.Done(): // base context canceled
		cancel()
		<-errC // wait for Run to return to avoid leaking the goroutine
		return noopStop, ctx.Err()
	case err := <-errC: // Run returned before becoming ready
		cancel()
		return noopStop, err
	case <-timeoutC: // startup timed out
		cancel()
		<-errC // wait for Run to return to avoid leaking the goroutine
		return noopStop, ErrStartupTimeout
	case <-initDone: // ready
	}

	stopFn := func() error {
		cancel()
		return <-errC
	}
	return stopFn, nil
}
