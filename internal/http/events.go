package http

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"sermoni/internal/events"
	"sermoni/internal/services"
	"strconv"

	"github.com/gorilla/mux"
)

const headerServiceToken = "Service-Token"

func getEvents(w http.ResponseWriter, r *http.Request) {
	// Create a mapping from service id to name
	/* Eventually?
	serviceIdName := make(map[int]string)
	services := services.GetAll()
	for _, service := range {
		serviceIdName[service.ID] = service.Name
	}
	*/

	events := events.GetAll()
	var s *services.Service
	for _, e := range events {
		// Could optimize by writing a single db.View, but this is more readable, so meh..
		s = services.GetByID(e.Service)
		e.ServiceName = s.Name
	}
	b, _ := json.Marshal(events)
	w.Write(b)
}

// TODO: This is still a placeholder!
func deleteEvent(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, _ := strconv.ParseUint(vars["id"], 10, 64)
	err := events.Delete(id)
	if err != nil {
		// TODO: Non-existing error is not an internal server error
		fmt.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
	}
}

func reportEvent(w http.ResponseWriter, r *http.Request) {
	tokens := r.Header[headerServiceToken]
	if len(tokens) == 0 {
		w.WriteHeader(http.StatusUnauthorized)
		msg := fmt.Sprintf("%v: No service token given\n", http.StatusUnauthorized)
		w.Write([]byte(msg))
		return
	}
	service := services.GetByToken(tokens[0])
	if service == nil {
		msg := fmt.Sprintf("%v: No service for the given token\n", http.StatusUnauthorized)
		w.Write([]byte(msg))
		return
	}
	content, err := ioutil.ReadAll(r.Body)
	check(err)
	// TODO: This should deffo not panic on error!
	event := new(events.Event)
	err = json.Unmarshal(content, event)
	check(err)
	event.Service = service.ID
	err = events.Add(event)
	check(err)
	log.Printf("New event registered, id = %v\n", event.ID)
}

// For later reference:
// Instead of anonymous structs for each JSON object to be written,
// a simple map can be used instead, something like (untested!):
//	b, _ := json.Marshal(&map[string]string{
//	 	"message": http.StatusUnauthorized + ": No valid token"
//	})
