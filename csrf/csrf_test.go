package csrf

import (
	"github.com/plimble/ace"
	"github.com/plimble/sessions/store/cookie"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestCSRFForm(t *testing.T) {
	assert := assert.New(t)

	token := ""

	a := ace.New()
	a.Session(cookie.NewCookieStore(), nil)
	CSRF(nil)

	a.GET("/", func(c *ace.C) {
		token = Token(c)
		c.JSON(200, nil)
	})

	a.POST("/", Validate, func(c *ace.C) {
		c.String(200, "passed")
	})

	r, _ := http.NewRequest("GET", "/", nil)
	w := httptest.NewRecorder()
	a.ServeHTTP(w, r)
	assert.NotEmpty(token)

	cookie := w.Header().Get("Set-Cookie")
	r, _ = http.NewRequest("POST", "/", nil)
	r.Header.Set("Cookie", cookie)
	r.ParseForm()
	r.PostForm.Set("csrf_token", token)
	w = httptest.NewRecorder()
	a.ServeHTTP(w, r)
	assert.Equal(200, w.Code)
	assert.Equal("passed", w.Body.String())

	cookie = w.Header().Get("Set-Cookie")
	r, _ = http.NewRequest("POST", "/", nil)
	r.Header.Set("Cookie", cookie)
	r.ParseForm()
	r.PostForm.Set("csrf_token", token)
	w = httptest.NewRecorder()
	a.ServeHTTP(w, r)
	assert.Equal(500, w.Code)
	assert.Equal("Invalid CSRF Token", w.Body.String())
}
