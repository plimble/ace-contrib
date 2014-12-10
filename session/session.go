package session

import (
	"github.com/gorilla/context"
	"github.com/plimble/ace"
)

func Sessions() ace.HandlerFunc {
	return func(c *ace.C) {
		defer context.Clear(c.Request)
		c.Next()
	}
}
