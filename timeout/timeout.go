package timeout

import (
	"time"
)

type Timeout struct {
	expire time.Time
}

func TimerEventSet(duration time.Duration) Timeout {
	return Timeout{expire: time.Now().Add(duration)}
}

// Devuelve true si el timeout ya expiró
func (t Timeout) TimerEventHasExpired() bool {
	return time.Now().After(t.expire)
}
