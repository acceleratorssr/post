package distributed_lock

import (
	"time"
)

type MutexOption func(*Mutex)

type Redsync struct {
	clients []*Client
}

func New(pools []*Client) *Redsync {
	return &Redsync{
		clients: pools,
	}
}

func (r *Redsync) NewMutex(name string, opts ...MutexOption) *Mutex {
	m := &Mutex{
		name:       name,
		expiration: 8 * time.Second,
		reTry: &FixIntervalRetry{
			Interval: 500 * time.Millisecond,
			Max:      32,
		},
		quorum:       len(r.clients)/2 + 1,
		genValueFunc: genValue,
		factor:       0.01,

		clients: r.clients,
	}
	for _, opt := range opts {
		opt(m)
	}
	return m
}

func WithExpiration(expiration time.Duration) MutexOption {
	return func(m *Mutex) {
		m.expiration = expiration
	}
}

func WithRetryStrategy(retry RetryStrategy) MutexOption {
	return func(m *Mutex) {
		m.reTry = retry
	}
}

func WithValueGenerator(genValueFunc func() (string, error)) MutexOption {
	return func(m *Mutex) {
		m.genValueFunc = genValueFunc
	}
}

func WithDriftFactor(factor float64) MutexOption {
	return func(m *Mutex) {
		m.factor = factor
	}
}
