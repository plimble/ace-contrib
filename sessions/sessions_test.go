package sessions_test

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/plimble/ace/sessions"
	"strings"

	"net/http"
	"net/http/httptest"
	"testing"
)

func init() {
	gin.SetMode(gin.TestMode)
}

func Test_Sessions(t *testing.T) {
	n := gin.New()

	store := sessions.NewCookieStore([]byte("secret123"))
	n.Use(sessions.Sessions("my_session", store))

	n.GET("/testsession", func(c *gin.Context) {
		session := sessions.GetSession(c.Request)
		session.Set("hello", "world")
		session.Session().Save(c.Request, c.Writer)
		fmt.Fprintf(c.Writer, "OK")
	})

	n.GET("/show", func(c *gin.Context) {
		session := sessions.GetSession(c.Request)
		if session.Get("hello") != "world" {
			t.Error("Session writing failed")
		}
	})

	res := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/testsession", nil)
	n.ServeHTTP(res, req)

	res2 := httptest.NewRecorder()
	req2, _ := http.NewRequest("GET", "/show", nil)
	req2.Header.Set("Cookie", res.Header().Get("Set-Cookie"))
	n.ServeHTTP(res2, req2)
}

func Test_SessionsDeleteValue(t *testing.T) {
	n := gin.New()

	store := sessions.NewCookieStore([]byte("secret123"))
	n.Use(sessions.Sessions("my_session", store))

	n.GET("/testsession", func(c *gin.Context) {
		session := sessions.GetSession(c.Request)
		session.Set("hello", "world")
		session.Delete("hello")
		session.Session().Save(c.Request, c.Writer)
		fmt.Fprintf(c.Writer, "OK")
	})

	n.GET("/show", func(c *gin.Context) {
		session := sessions.GetSession(c.Request)
		if session.Get("hello") == "world" {
			t.Error("Session value deleting failed")
		}
		fmt.Fprintf(c.Writer, "OK")
	})

	res := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/testsession", nil)
	n.ServeHTTP(res, req)

	res2 := httptest.NewRecorder()
	req2, _ := http.NewRequest("GET", "/show", nil)
	req2.Header.Set("Cookie", res.Header().Get("Set-Cookie"))
	n.ServeHTTP(res2, req2)
}

func Test_Options(t *testing.T) {
	n := gin.New()
	store := sessions.NewCookieStore([]byte("secret123"))
	store.Options(sessions.Options{
		Domain: "negroni-sessions.goincremental.com",
	})

	n.Use(sessions.Sessions("my_session", store))

	n.GET("/", func(c *gin.Context) {
		session := sessions.GetSession(c.Request)
		session.Set("hello", "world")
		session.Options(sessions.Options{
			Path: "/foo/bar/bat",
		})
		session.Session().Save(c.Request, c.Writer)
		fmt.Fprintf(c.Writer, "OK")
	})

	n.GET("/foo", func(c *gin.Context) {
		session := sessions.GetSession(c.Request)
		session.Set("hello", "world")
		session.Session().Save(c.Request, c.Writer)
		fmt.Fprintf(c.Writer, "OK")
	})

	res := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/", nil)
	n.ServeHTTP(res, req)

	res2 := httptest.NewRecorder()
	req2, _ := http.NewRequest("GET", "/foo", nil)
	n.ServeHTTP(res2, req2)

	s := strings.Split(res.Header().Get("Set-Cookie"), ";")

	if s[1] != " Path=/foo/bar/bat" {
		t.Error("Error writing path with options:", s[1])
	}

	s = strings.Split(res2.Header().Get("Set-Cookie"), ";")
	if s[1] != " Domain=negroni-sessions.goincremental.com" {
		t.Error("Error writing domain with options:", s[1])
	}
}

func Test_Flashes(t *testing.T) {
	n := gin.New()

	store := sessions.NewCookieStore([]byte("secret123"))
	n.Use(sessions.Sessions("my_session", store))

	n.GET("/set", func(c *gin.Context) {
		session := sessions.GetSession(c.Request)
		session.AddFlash("hello world")
		session.Session().Save(c.Request, c.Writer)
		fmt.Fprintf(c.Writer, "OK")
	})

	n.GET("/show", func(c *gin.Context) {
		session := sessions.GetSession(c.Request)
		l := len(session.Flashes())
		if l != 1 {
			t.Error("Flashes count does not equal 1. Equals ", l)
		}
		fmt.Fprintf(c.Writer, "OK")
	})

	n.GET("/showagain", func(c *gin.Context) {
		session := sessions.GetSession(c.Request)
		l := len(session.Flashes())
		if l != 0 {
			t.Error("flashes count is not 0 after reading. Equals ", l)
		}
		fmt.Fprintf(c.Writer, "OK")
	})

	res := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/set", nil)
	n.ServeHTTP(res, req)

	res2 := httptest.NewRecorder()
	req2, _ := http.NewRequest("GET", "/show", nil)
	req2.Header.Set("Cookie", res.Header().Get("Set-Cookie"))
	n.ServeHTTP(res2, req2)

	res3 := httptest.NewRecorder()
	req3, _ := http.NewRequest("GET", "/showagain", nil)
	req3.Header.Set("Cookie", res2.Header().Get("Set-Cookie"))
	n.ServeHTTP(res3, req3)
}

func Test_SessionsClear(t *testing.T) {
	n := gin.New()
	data := map[string]string{
		"hello":  "world",
		"foo":    "bar",
		"apples": "oranges",
	}

	store := sessions.NewCookieStore([]byte("secret123"))
	n.Use(sessions.Sessions("my_session", store))

	n.GET("/testsession", func(c *gin.Context) {
		session := sessions.GetSession(c.Request)
		for k, v := range data {
			session.Set(k, v)
		}
		session.Clear()
		session.Session().Save(c.Request, c.Writer)
		fmt.Fprintf(c.Writer, "OK")
	})

	n.GET("/show", func(c *gin.Context) {
		session := sessions.GetSession(c.Request)
		for k, v := range data {
			if session.Get(k) == v {
				t.Fatal("Session clear failed")
			}
		}
		fmt.Fprintf(c.Writer, "OK")
	})

	res := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/testsession", nil)
	n.ServeHTTP(res, req)

	res2 := httptest.NewRecorder()
	req2, _ := http.NewRequest("GET", "/show", nil)
	req2.Header.Set("Cookie", res.Header().Get("Set-Cookie"))
	n.ServeHTTP(res2, req2)
}
