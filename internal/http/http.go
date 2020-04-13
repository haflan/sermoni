package http

import (
	"fmt"
	"log"
	"net/http"

	"sermoni/internal/config"

	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
)

var conf *config.Config
var store *sessions.CookieStore

// StartServer initializes the session store given the session key and starts
// the server at the given port
func StartServer(port int) {
	conf = config.GetConfig()
	store = sessions.NewCookieStore(conf.SessionKey)

	router := mux.NewRouter()
	router.HandleFunc("/", homeHandler)
	router.HandleFunc("/init", initHandler).Methods(http.MethodGet)
	router.HandleFunc("/login", loginHandler).Methods(http.MethodPost)
	router.Handle("/logout", auth(logoutHandler))

	router.Handle("/services", auth(getServices)).Methods(http.MethodGet)
	router.Handle("/services", auth(postService)).Methods(http.MethodPost)
	router.Handle("/services/{id:[0-9]+}", auth(deleteService)).Methods(http.MethodDelete)
	//router.Handle("/services/{id:[0-9]+}", putService).Methods("PUT") (TODO)

	router.Handle("/events", auth(getEvents)).Methods(http.MethodGet)
	router.Handle("/events/{id:[0-9]+}", auth(deleteEvent)).Methods(http.MethodDelete)
	// POSTS to /events is how services should report events to the monitor
	// This should not be accessible from the website
	router.HandleFunc("/events", reportEvent).Methods(http.MethodPost)

	http.Handle("/", router)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%v", port), nil))
}

func homeHandler(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	w.Write(getWebsite())
	return
}

func check(err error) {
	if err != nil {
		panic(err)
	}
}
