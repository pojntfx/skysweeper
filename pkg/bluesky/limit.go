package bluesky

import (
	"context"
	"sync"
	"time"
)

type Limiter struct {
	ctx context.Context

	globalLimit         int
	availablePoints     int
	availablePointsLock *sync.Cond

	ticker *time.Ticker
}

func NewLimiter(
	ctx context.Context,
	globalLimit int,
	resetInterval time.Duration,
) *Limiter {
	return &Limiter{
		ctx: ctx,

		globalLimit:         globalLimit,
		availablePoints:     globalLimit,
		availablePointsLock: &sync.Cond{},

		ticker: time.NewTicker(resetInterval),
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

func (l *Limiter) Spend(points int) bool {
	l.availablePointsLock.L.Lock()
	if l.availablePoints-points <= 0 {
		l.availablePointsLock.Wait()
	}

	if l.availablePoints < 0 {
		return true // Context cancelled
	}

	l.availablePoints = l.availablePoints - points

	l.availablePointsLock.L.Unlock()

	return false
}
