package config

import (
	"net/http"

	"github.com/gorilla/sessions"
)

type CookieSessionStore struct {
	Store *sessions.CookieStore
}

func (s *CookieSessionStore) Get(r *http.Request, name string) (*sessions.Session, error) {
	return s.Store.Get(r, name)
}

func (s *CookieSessionStore) Save(r *http.Request, w http.ResponseWriter, session *sessions.Session) error {
	return s.Store.Save(r, w, session)
}

var Sessions = sessions.NewCookieStore([]byte("secret-key")) // secret-key should be hidden in a production environment.

func init() {
	Sessions.Options = &sessions.Options{
		Path:     "/",
		MaxAge:   2592000,
		Secure:   false,
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
	}
}
