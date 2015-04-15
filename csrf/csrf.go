package csrf

import (
	"encoding/base64"
	"github.com/plimble/ace"
	"time"
)

var csrfOptions *CSRFOptions

type CSRFOptions struct {
	// Name used for create the session. Default csrf_token
	SessionName string
	// If true, send token via X-CSRFToken header. Default false
	SetHeader bool
	// HTTP header used to set and get token. Default X-CSRFToken
	Header string
	// Form value used to set and get token. Default csrf_token
	Form string
	//ErrorHandler when csrf is invalid
	ErrorHandler ace.HandlerFunc
}

func CSRF(options *CSRFOptions) {
	csrfOptions = options

	if csrfOptions == nil {
		csrfOptions = &CSRFOptions{}
	}

	if csrfOptions.SessionName == "" {
		csrfOptions.SessionName = "csrf_token"
	}

	if csrfOptions.Header == "" {
		csrfOptions.Header = "X-CSRFToken"
	}

	if csrfOptions.Form == "" {
		csrfOptions.Form = "csrf_token"
	}

	if csrfOptions.ErrorHandler == nil {
		csrfOptions.ErrorHandler = func(c *ace.C) {
			c.String(500, "Invalid CSRF Token")
		}
	}
}

func Token(c *ace.C) string {
	session, _ := c.Sessions.Get(csrfOptions.SessionName)

	token := generateNewToken(session.ID)
	session.Set("token", token)
	session.Set("id", session.ID)

	if csrfOptions.SetHeader {
		c.Writer.Header().Set(csrfOptions.Header, token)
	}

	return token
}

func Validate(c *ace.C) {
	session, _ := c.Sessions.Get(csrfOptions.SessionName)

	if token := c.Request.Header.Get(csrfOptions.Header); token != "" {
		if validateCSRFToken(session.GetString("id", ""), session.GetString("token", ""), token) {
			c.Next()
			return
		}
	} else if token := c.MustPostString(csrfOptions.Form, ""); token != "" {
		if validateCSRFToken(session.GetString("id", ""), session.GetString("token", ""), token) {
			c.Next()
			return
		}
	}

	csrfOptions.ErrorHandler(c)
	c.Abort()
}

func validateCSRFToken(sessionID, sessionToken, token string) bool {
	if sessionToken != token {
		return false
	}

	b, err := base64.StdEncoding.DecodeString(token)
	if err != nil {
		return false
	}

	id := ""
	timeStr := ""
	for i := 0; i < len(b); i++ {
		if b[i] == '#' {
			id = string(b[:i])
			timeStr = string(b[i+1:])
			break
		}
	}

	if id != sessionID {
		return false
	}

	t, err := time.Parse(time.RFC3339, timeStr)
	if err != nil {
		return false
	}

	return t.UTC().After(time.Now().UTC())
}

func generateNewToken(id string) string {
	pool := ace.GetPool()
	buf := pool.Get()
	defer pool.Put(buf)

	buf.WriteString(id)
	buf.WriteString("#")
	buf.WriteString(time.Now().UTC().Add(time.Hour * 24).Format(time.RFC3339))

	return base64.StdEncoding.EncodeToString(buf.Bytes())
}
