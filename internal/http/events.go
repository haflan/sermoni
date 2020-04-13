package http

import (
	"encoding/json"
	"net/http"
	"sermoni/internal/events"
)

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
	json.Marshal(events)
	return
}

// TODO: This is still a placeholder!
func deleteEvent(w http.ResponseWriter, r *http.Request) {
	/*
		vars := mux.Vars(r)
		id, _ := strconv.ParseUint(vars["id"], 10, 64)
		err := events.Delete(id)
		fmt.Println(err)
	*/
	return
}

func reportEvent(w http.ResponseWriter, r *http.Request) {
	return
}
