package config

import (
	"net/http"

	"github.com/gorilla/sessions"
)

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
