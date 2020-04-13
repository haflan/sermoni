package http

import (
	"crypto/sha256"
	"encoding/json"
	"io/ioutil"
	"math/rand"
	"net/http"
	"strings"
	"time"

	"github.com/gorilla/sessions"
)

// Deal with login, logout, and general security stuff

const (
	keyAuthenticated = "authenticated"
	keyCSRFToken     = "csrfToken"
	keyPassphrase    = "passphrase"
	keySessionName   = "session"
	headerCSRFToken  = "X-Csrf-Token"
)

// initHandler checks two things:
// 1. If a CSRF token exists for the given session. Otherwise it creates it
// 2. Whether the session is authenticated
// It then returns an object on the form {"auth": true, "csrftoken": "<long string>"}
// This is requested immediately when the website is loaded.
func initHandler(w http.ResponseWriter, r *http.Request) {
	session, _ := store.Get(r, "session")
	val := session.Values[keyCSRFToken]
	token, ok := val.(string)
	if !ok {
		token = temporary32CharRandomString()
		session.Values[keyCSRFToken] = token
		session.Save(r, w) // TODO: Error handling, as always
	}
	b, _ := json.Marshal(struct {
		CSRFToken     string `json:"csrftoken"`
		Authenticated bool   `json:"authenticated"`
	}{
		token,
		authorized(session),
	})
	w.Write(b)
}

func loginHandler(w http.ResponseWriter, r *http.Request) {
	session, _ := store.Get(r, keySessionName)
	if authorized(session) {
		return
	}
	defer r.Body.Close() // needed?
	content, err := ioutil.ReadAll(r.Body)
	check(err)
	var data map[string]string
	json.Unmarshal(content, &data)
	passphrase := data[keyPassphrase]
	passhash := sha256.Sum256([]byte(passphrase))
	if string(passhash[:]) == string(conf.PassHash) {
		session.Values[keyAuthenticated] = true
		err = session.Save(r, w)
		check(err)
		w.WriteHeader(http.StatusOK) // Not needed, just for readability?
	} else {
		w.WriteHeader(http.StatusUnauthorized)
	}
}

func logoutHandler(w http.ResponseWriter, r *http.Request) {
	session, _ := store.Get(r, keySessionName)
	session.Values[keyAuthenticated] = false
	err := session.Save(r, w)
	check(err)
	w.Write([]byte("Logged out"))
}

func authorized(session *sessions.Session) bool {
	val := session.Values[keyAuthenticated]
	auth, ok := val.(bool)
	return ok && auth
}

// Middleware for the simple sermoni authentication scheme
func auth(handler http.HandlerFunc) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// store is the global CookieStore
		session, _ := store.Get(r, keySessionName)
		if !authorized(session) {
			status := http.StatusUnauthorized
			http.Error(w, http.StatusText(status), status)
			return
		}
		handler.ServeHTTP(w, r)
	})
}

func csrfCheckPassed(r *http.Request, session *sessions.Session) bool {
	// CSRF protect anything but GET requests
	if r.Method == http.MethodGet {
		return true
	}
	val := session.Values[keyCSRFToken]
	rightToken, ok := val.(string)
	if !ok {
		panic("no CSRF token found")
	}
	if tokenHeader := r.Header[headerCSRFToken]; tokenHeader == nil {
		return false
	} else {
		return tokenHeader[0] == rightToken
	}
}

// not cryptosecure, only for testing!
// thanks: https://yourbasic.org/golang/generate-random-string/
func temporary32CharRandomString() string {
	rand.Seed(time.Now().UnixNano())
	chars := []rune("ABCDEFGHIJKLMNOPQRSTUVWXYZ" +
		"abcdefghijklmnopqrstuvwxyz" +
		"0123456789")
	length := 32
	var b strings.Builder
	for i := 0; i < length; i++ {
		b.WriteRune(chars[rand.Intn(len(chars))])
	}
	return b.String()
}