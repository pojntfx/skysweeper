package bluesky

import (
	"context"
	"sync"
	"time"
)

const (
	PointsCreate = 3
	PointsUpdate = 2
	PointsDelete = 1

	PointsGet = 1 // Technically not specified, so its assumed that its the equivalent of delete (see https://atproto.com/blog/rate-limits-pds-v3)
)

type Limiter struct {
	ctx context.Context

	globalLimit         int
	availablePoints     int
	availablePointsLock *sync.Cond

	ticker *time.Ticker

	onWaitingForReset func() error
}

func NewLimiter(
	ctx context.Context,

	globalLimit int,
	resetInterval time.Duration,

	onWaitingForReset func() error,
) *Limiter {
	return &Limiter{
		ctx: ctx,

		globalLimit: globalLimit,

		availablePoints:     globalLimit,
		availablePointsLock: sync.NewCond(&sync.Mutex{}),

		ticker: time.NewTicker(resetInterval),

		onWaitingForReset: onWaitingForReset,
	}
}

func (l *Limiter) Open() {
	for {
		select {
		case <-l.ctx.Done():
			l.ticker.Stop()

			l.availablePointsLock.L.Lock()
			l.availablePoints = -1
			l.availablePointsLock.Broadcast()
			l.availablePointsLock.L.Unlock()

			return
		case <-l.ticker.C:
			l.availablePointsLock.L.Lock()
			l.availablePoints = l.globalLimit
			l.availablePointsLock.Broadcast()
			l.availablePointsLock.L.Unlock()
		}
	}
}

func (l *Limiter) Spend(points int) error {
	l.availablePointsLock.L.Lock()
	if l.availablePoints-points <= 0 {
		if l.onWaitingForReset != nil {
			if err := l.onWaitingForReset(); err != nil {
				return err
			}
		}

		l.availablePointsLock.Wait()
	}

	if l.availablePoints < 0 {
		return context.Canceled // Context cancelled
	}

	l.availablePoints = l.availablePoints - points

	l.availablePointsLock.L.Unlock()

	return nil
}
