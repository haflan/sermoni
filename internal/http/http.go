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

	auth := AuthHandler(store)
	router := mux.NewRouter()
	router.HandleFunc("/", homeHandler)
	router.HandleFunc("/login", loginHandler)
	//router.HandleFunc("/logout", logoutHandler)
	router.Handle("/logout", auth(http.HandlerFunc(logoutHandler)))

	router.Handle("/services", auth(http.HandlerFunc(getServices))).Methods("GET")
	router.Handle("/services", auth(http.HandlerFunc(postService))).Methods("POST")
	router.HandleFunc("/services/{id:[0-9]+}", deleteService).Methods("DELETE")
	//router.HandleFunc("/services/{id:[0-9]+}", putService).Methods("PUT") (TODO)

	router.HandleFunc("/events", getEvents).Methods("GET")
	router.HandleFunc("/events/{id:[0-9]+}", deleteEvent).Methods("DELETE")

	router.HandleFunc("/report", reportEvent).Methods("POST")
	http.Handle("/", router)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%v", port), nil))
}

func homeHandler(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	w.Write(getWebsite())
	return
}
