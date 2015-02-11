package csrf

import (
	"github.com/plimble/ace"
	"github.com/plimble/csrf"
)

type CSRFOptions struct {
	FailedHandler ace.HandlerFunc
}

func CSRF(options *CSRFOptions) ace.HandlerFunc {
	if options.FailedHandler == nil {
		options.FailedHandler = defaultCSRFFailedHandler
	}

	cs := csrf.New()

	return func(c *ace.C) {
		defer csrf.ClearContext(c.Request)

		if !cs.Check(c.Writer, c.Request) {
			options.FailedHandler(c)
			c.Abort()
		}

		c.Next()
	}
}

func defaultCSRFFailedHandler(c *ace.C) {
	c.String(500, "Invalid CSRF Token")
}
