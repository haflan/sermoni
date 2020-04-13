package http

import (
	"encoding/json"
	"net/http"
	"sermoni/internal/services"
	"strconv"

	"github.com/gorilla/mux"
)

func getServices(w http.ResponseWriter, r *http.Request) {
	s := services.GetAll()
	data, err := json.Marshal(s)
	check(err)
	w.Write(data)
}
func postService(w http.ResponseWriter, r *http.Request) {

}
func deleteService(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, _ := strconv.ParseUint(vars["id"], 10, 64)
	err := services.Delete(id)
	if err != nil {
		// TODO: Non-existing error is not an internal server error
		w.WriteHeader(http.StatusInternalServerError)
	}
}
