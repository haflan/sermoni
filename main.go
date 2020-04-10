package main

import (
	"flag"
	"fmt"
	"log"
	"sermoni/database"
	"sermoni/services"
)

func main() {
	// TODO: Use getopt package instead of flags?
	port := flag.Int("p", 8080, "Port")
	dbFile := flag.String("d", "sermoni.db", "Database file")
	//password := flag.String("w", "", "Password for the web interface")
	flag.Parse()
	if err := database.Init(*dbFile); err != nil {
		log.Fatal(err)
	}
	defer database.Close()
	fmt.Printf("Server running on port %v\n", *port)
	services.Add("testing", services.Service{
		Name:        "Test name",
		Description: "This is the description, yay",
	})
	testService := services.Get("testing")
	fmt.Printf("test service: %+v\n", testService)
}
