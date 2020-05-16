package http

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

// Deal with login, logout, and general security stuff

const (
	keyAuthenticated = "authenticated"
	keyCSRFToken     = "csrfToken"
	keyPassphrase    = "passphrase"
	keySessionName   = "session"
	headerCSRFToken  = "X-Csrf-Token"
	headerPassHash   = "Pass-Hash"
)

// initHandler checks two things:
// 1. If a CSRF token exists for the given session. Otherwise it creates it
// 2. Whether the session is authenticated
// It then returns an object on the form {"auth": true, "csrftoken": "<long string>", "pagetitle": "Page title"}
// This is requested immediately when the website is loaded.
func initHandler(w http.ResponseWriter, r *http.Request) {
	session, _ := store.Get(r, keySessionName)
	val := session.Values[keyCSRFToken]
	token, ok := val.(string)
	if !ok {
		token = generateCSRFToken()
		session.Values[keyCSRFToken] = token
		session.Save(r, w) // TODO: Error handling, as always
	}
	b, _ := json.Marshal(struct {
		PageTitle     string `json:"pagetitle"`
		CSRFToken     string `json:"csrftoken"`
		Authenticated bool   `json:"authenticated"`
	}{
		string(conf.PageTitle),
		token,
		authorized(r),
	})
	w.Write(b)
}

func loginHandler(w http.ResponseWriter, r *http.Request) {
	if authorized(r) {
		return
	}
	defer r.Body.Close() // needed?
	content, err := ioutil.ReadAll(r.Body)
	check(err)
	var data map[string]string
	json.Unmarshal(content, &data)
	passphrase := data[keyPassphrase]
	passhash := sha256.Sum256([]byte(passphrase))
	// Converting byte slices to string is fine, it's how bytes.Equal does it
	if string(passhash[:]) == string(conf.PassHash) {
		session, _ := store.Get(r, keySessionName)
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
	b, _ := json.Marshal(struct {
		Info string `json:"info"`
	}{"Logged out"})
	w.Write(b)
}

func authorized(r *http.Request) bool {
	session, _ := store.Get(r, keySessionName)
	val := session.Values[keyAuthenticated]
	auth, ok := val.(bool)
	return ok && auth
}

// Middleware for the simple sermoni authentication scheme
// Accepts authentication both with session cookie (typically web browser) and
// a header on the format `Pass-Hash: <sha256sum of pass phrase>`
func auth(handler http.HandlerFunc) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		passHashHeader := r.Header[headerPassHash]
		if passHashHeader != nil {
			// compare to hex formatted version of correct hash
			if passHashHeader[0] != fmt.Sprintf("%x", conf.PassHash) {
				status := http.StatusUnauthorized
				http.Error(w, "Invalid passphrase hash", status)
				return
			}
		} else {
			// store is the global CookieStore
			if !authorized(r) || !csrfCheckPassed(r) {
				status := http.StatusUnauthorized
				http.Error(w, http.StatusText(status), status)
				return
			}
		}
		handler.ServeHTTP(w, r)
	})
}

func csrfCheckPassed(r *http.Request) bool {
	session, _ := store.Get(r, keySessionName)
	// CSRF protect anything but GET requests
	if r.Method == http.MethodGet {
		return true
	}
	val := session.Values[keyCSRFToken]
	rightToken, ok := val.(string)
	if !ok {
		panic("no CSRF token found")
	}
	tokenHeader := r.Header[headerCSRFToken]
	if tokenHeader == nil {
		return false
	}
	return tokenHeader[0] == rightToken
}

func generateCSRFToken() string {
	randBytes := make([]byte, 32)
	rand.Read(randBytes)
	return hex.EncodeToString(randBytes)
}
