package http

import (
	"encoding/json"
	"log"
	"net/http"
	"math/rand"
	"strings"
	"time"

	"github.com/gorilla/sessions"
)

// Deal with login, logout, and general security stuff

// initHandler checks two things: 
// 1. If a CSRF token exists for the given session. Otherwise it creates it
// 2. Whether the session is authenticated
// It then returns an object on the form {"auth": true, "csrftoken": "<long string>"}
// This is requested immediately when the website is loaded.
func initHandler(w http.ResponseWriter, r *http.Request) {
	session, _ := store.Get(r, "session")
	val := session.Values["csrftoken"]
	token, ok := val.(string)
	if !ok {
		token = temporary32CharRandomString()
		session.Values["csrftoken"] = token
		session.Save(r, w) // TODO: Error handling, as always
	}
	b, _ := json.Marshal(struct {
		CSRFToken string `json:"csrftoken"`
		Authenticated bool `json:"authenticated"`
	}{
		token,
		authorized(session),
	})
	w.Write(b)
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

func authorized(session *sessions.Session) bool {
	val := session.Values["authenticated"]
	auth, ok := val.(bool)
	return ok && auth
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

// not cryptosecure, only for testing!
// thanks: https://yourbasic.org/golang/generate-random-string/
func temporary32CharRandomString() string {
	rand.Seed(time.Now().UnixNano())
	chars := []rune("ABCDEFGHIJKLMNOPQRSTUVWXYZÅÄÖ" +
	    "abcdefghijklmnopqrstuvwxyzåäö" +
	    "0123456789")
	length := 32
	var b strings.Builder
	for i := 0; i < length; i++ {
	    b.WriteRune(chars[rand.Intn(len(chars))])
	}
	return b.String()
}
