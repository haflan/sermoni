package http

import (
	"encoding/json"
	"io/ioutil"
	"log"
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
	content, err := ioutil.ReadAll(r.Body)
	check(err)
	service := new(services.Service)
	// TODO: Handle json parse error
	err = json.Unmarshal(content, service)
	check(err)
	services.Add(service)
	w.WriteHeader(http.StatusCreated)
}

func deleteService(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, _ := strconv.ParseUint(vars["id"], 10, 64)
	err := services.Delete(id)
	if err != nil {
		// TODO: "Non-existing" error is not an internal server error
		log.Printf("deleteService error: %v\n", err)
		w.WriteHeader(http.StatusInternalServerError)
	}
}
