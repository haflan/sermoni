package http

import (
	"net/http"
)

func getEvents(w http.ResponseWriter, r *http.Request) {
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
