package cors

import (
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/plimble/ace"
)

var (
	defaultAllowHeaders = []string{"Origin", "Accept", "Content-Type", "Authorization"}
	defaultAllowMethods = []string{"GET", "POST", "PUT", "DELETE", "PATCH", "HEAD"}
)

// Options stores configurations
type Options struct {
	AllowOrigins     []string
	AllowCredentials bool
	AllowMethods     []string
	AllowHeaders     []string
	ExposeHeaders    []string
	MaxAge           time.Duration
}

// Middleware sets CORS headers for every request
func Cors(options Options) ace.HandlerFunc {
	if options.AllowHeaders == nil {
		options.AllowHeaders = defaultAllowHeaders
	}

	if options.AllowMethods == nil || len(options.AllowMethods) == 0 {
		options.AllowMethods = defaultAllowMethods
	}

	allowOrigin := "*"
	if len(options.AllowOrigins) > 0 {
		allowOrigin = strings.Join(options.AllowOrigins, " ")
	}
	exposeHeader := strings.Join(options.ExposeHeaders, ",")
	allowMethod := strings.Join(options.AllowMethods, ",")
	allowHeader := strings.Join(options.AllowHeaders, ",")

	maxAge := ""
	if options.MaxAge > time.Duration(0) {
		maxAge = strconv.FormatInt(int64(options.MaxAge/time.Second), 10)
	}

	return func(c *ace.C) {
		// origin := c.Request.Header.Get("Origin")
		// requestMethod := c.Request.Header.Get("Access-Control-Request-Method")
		// requestHeaders := c.Request.Header.Get("Access-Control-Request-Headers")
		c.Writer.Header().Set("Access-Control-Allow-Origin", allowOrigin)

		if options.AllowCredentials {
			c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		}

		if exposeHeader != "" {
			c.Writer.Header().Set("Access-Control-Expose-Headers", exposeHeader)
		}

		if c.Request.Method == "OPTIONS" {
			c.Writer.Header().Set("Access-Control-Allow-Methods", allowMethod)
			c.Writer.Header().Set("Access-Control-Allow-Headers", allowHeader)

			if maxAge != "" {
				c.Writer.Header().Set("Access-Control-Max-Age", maxAge)
			}

			c.Abort(http.StatusOK)
		} else {
			c.Next()
		}
	}
}
