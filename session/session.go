package session

import (
	"github.com/gorilla/context"
	"github.com/plimble/copter"
)

func Sessions() copter.HandlerFunc {
	return func(c *copter.C) {
		defer context.Clear(c.Request)
		c.Next()
	}
}
