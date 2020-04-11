package main

import (
	"flag"
	"fmt"
	"log"
	"sermoni/internal/database"
	smhttp "sermoni/internal/http"
)

var (
	port   = flag.Int("p", 8080, "Port")
	dbFile = flag.String("d", "sermoni.db", "Database file")
)

func main() {
	// TODO: Use getopt package instead of flags?
	//password := flag.String("w", "", "Password for the web interface")
	flag.Parse()
	defer database.Close()
	if err := database.Init(*dbFile); err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Server running on port %v\n", *port)
	smhttp.StartServer(*port)
}
