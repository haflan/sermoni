package http

import (
	"log"
	"net/http"

	"github.com/gorilla/sessions"
)

func authorized(session *sessions.Session) bool {
	val := session.Values["authenticated"]
	auth, ok := val.(bool)
	return ok && auth
}

func loginHandler(w http.ResponseWriter, r *http.Request) {
	session, _ := store.Get(r, "session")
	if authorized(session) {
		log.Println("Authenticated session requested website")
		w.Write([]byte("logged in"))
	} else {
		log.Println("New session requested website")
		session.Values["authenticated"] = true
		log.Println(session.Save(r, w))
		w.Write([]byte("Not logged in"))
	}
	return
}

func logoutHandler(w http.ResponseWriter, r *http.Request) {
	session, _ := store.Get(r, "session")
	session.Values["authenticated"] = false
	log.Println(session.Save(r, w))
	w.Write([]byte("Logged out"))
}

// Middleware for the simple sermoni authentication scheme
func auth(handler http.HandlerFunc) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Store is the global CookieStore
		session, _ := store.Get(r, "session")
		if !authorized(session) {
			status := http.StatusUnauthorized
			http.Error(w, http.StatusText(status), status)
			return
		}
		handler.ServeHTTP(w, r)
	})
}
