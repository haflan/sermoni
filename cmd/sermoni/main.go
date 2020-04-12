package main

import (
	"flag"
	"log"
	"sermoni/internal/database"
	smhttp "sermoni/internal/http"
)

// Flags
var (
	port   = flag.Int("p", 8080, "Port")
	dbFile = flag.String("d", "sermoni.db", "Database file")
)

func main() {
	// TODO: Use getopt package instead of flags?
	//password := flag.String("w", "", "Password for the web interface")
	flag.Parse()
	database.Init(*dbFile)
	defer database.Close()
	log.Printf("Server started listening on port %v\n", *port)
	smhttp.StartServer(*port)
}
