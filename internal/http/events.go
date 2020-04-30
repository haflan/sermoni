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
	"time"

	"github.com/gorilla/mux"
)

const headerServiceToken = "Service-Token"

func now() uint64 {
	return uint64(time.Now().UnixNano() / 1e6)
}

// getEvents fetches all events in the database and checks timestamp for the last report
// from each service. If the time since last report is more than the expected report interval,
// a "live" event is generated to inform that no report has been received
func getEvents(w http.ResponseWriter, r *http.Request) {
	// Create mappings from service id to the actual service and to the last service report
	serviceIDMap := make(map[uint64]*services.Service)
	serviceIsLate := make(map[uint64]bool)
	servs := services.GetAll()
	for _, service := range servs {
		serviceIDMap[service.ID] = service
		serviceIsLate[service.ID] = service.ExpectationPeriod != 0
	}

	var s *services.Service
	events := events.GetAll()
	for _, e := range events {
		sid := e.Service
		s = serviceIDMap[sid]
		e.ServiceName = s.Name
		if serviceIsLate[sid] && now()-e.Timestamp < s.ExpectationPeriod {
			serviceIsLate[sid] = false
		}
	}
	for sid, late := range serviceIsLate {
		if late {
			s := serviceIDMap[sid]
			events = append(events, generateLateEvent(s))
		}
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
	event.Timestamp = now()
	err = events.Add(event)
	check(err)
	log.Printf("New event registered, id = %v\n", event.ID)

	if service.MaxNumberEvents < 1 {
		return
	}
	// Find whether MaxNumberEvents is reached and delete the first event if so
	var numEvents, firstEventID uint64
	es := events.GetAll()
	for _, e := range es {
		if e.Service == service.ID {
			if firstEventID == 0 {
				firstEventID = e.ID
			}
			numEvents++
		}
	}
	if numEvents > service.MaxNumberEvents {
		events.Delete(firstEventID)
		log.Printf("MaxNumberEvents reached for service %v. Deleting first event", service.ID)
	}
}

func generateLateEvent(s *services.Service) *events.Event {
	return &events.Event{
		ID:        0,
		Service:   s.ID,
		Timestamp: now(),
		Status:    "late",
		Title:     "Expectation not met",
		Details: s.Name + " has failed to report within the expected internal." +
			"Something is probably wrong.",
		ServiceName: s.Name,
	}
}

// For later reference:
// Instead of anonymous structs for each JSON object to be written,
// a simple map can be used instead, something like (untested!):
//	b, _ := json.Marshal(&map[string]string{
//	 	"message": http.StatusUnauthorized + ": No valid token"
//	})
