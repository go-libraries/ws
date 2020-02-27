package ws

import "time"

func getPingPeriod(duration time.Duration) time.Duration {
	return (duration * 9) / 10
}
