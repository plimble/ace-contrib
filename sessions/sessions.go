package sessions

import (
	"github.com/boj/redistore"
	"github.com/gorilla/context"
	"github.com/gorilla/sessions"
	"github.com/plimble/ace"
	"net/http"
)

type Interface interface {
	Get(r *http.Request) *sessions.Session
}

type session struct {
	store sessions.Store
	name  string
}

func NewCookieSession(name string, key []byte) *session {
	return &session{sessions.NewCookieStore(key), name}
}

func NewFileSystemSession(name string, path string, key []byte) *session {
	return &session{sessions.NewFilesystemStore(path, key), name}
}

func NewRedisSession(name string, size int, network, address, password string, keyPairs []byte) *session {
	store, err := redistore.NewRediStore(size, network, address, password, keyPairs)
	if err != nil {
		panic(err)
	}
	return &session{store, name}
}

func Sessions() ace.HandlerFunc {
	return func(c *ace.C) {
		defer context.Clear(c.Request)
		c.Next()
	}
}

func (s *session) Get(r *http.Request) *sessions.Session {
	session, _ := s.store.Get(r, s.name)
	return session
}
