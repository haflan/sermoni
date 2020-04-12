package http

import (
	"fmt"
	"log"
	"net/http"
)

// StartServer starts the server at the given port
func StartServer(port int) {
	http.HandleFunc("/", staticHandler)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%v", port), nil))
}

func staticHandler(res http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	res.Write(getWebsite())
	return
}
