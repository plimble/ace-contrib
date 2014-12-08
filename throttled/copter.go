package throttled

import (
	"github.com/plimble/copter"
)

func Throttle(t *Throttler) copter.HandlerFunc {
	t.Start()
	return func(c *copter.C) {
		if t.Check(c.Writer, c.Request) {
			c.Next()
		}
	}
}
