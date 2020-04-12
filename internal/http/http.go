package http

import (
	"fmt"
	"log"
	"net/http"

	"sermoni/internal/config"

	"github.com/gorilla/sessions"
)

var store *sessions.CookieStore
var conf *config.Config

// StartServer initializes the session store given the session key and starts
// the server at the given port
func StartServer(port int) {
	conf = config.GetConfig()
	store = sessions.NewCookieStore(conf.SessionKey)

	http.HandleFunc("/", homeHandler)
	http.HandleFunc("/login", loginHandler)

	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%v", port), nil))
}

func homeHandler(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	w.Write(getWebsite())
	return
}

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
