package sugars

import (
	"context"
	"time"
)

// Pulse -
func Pulse(ctx context.Context, interval time.Duration) (func() bool, func()) {
	ctx, cancel := context.WithCancel(ctx)
	ticker := time.NewTicker(interval)

	pulser := func() bool {
		select {
		case <-ctx.Done():
			ticker.Stop()
			return false
		case _, ok := <-ticker.C:
			return ok
		}
	}

	stop := func() {
		ticker.Stop()
		cancel()
	}

	return pulser, stop
}
