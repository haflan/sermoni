package main

import (
	"flag"
	"log"
	"sermoni/internal/config"
	"sermoni/internal/database"
	smhttp "sermoni/internal/http"
)

// Flags
var (
	port   = flag.Int("p", 8080, "Port")
	dbFile = flag.String("d", "sermoni.db", "Database file")
	passphrase = flag.String("pass", "", "Set passphrase")
	pageTitle = flag.String("title", "", "Set page title")
)

func main() {
	// TODO: Use getopt package instead of flags?
	//password := flag.String("w", "", "Password for the web interface")
	flag.Parse()
	configured := database.Open(*dbFile)
	if !configured {
		log.Printf("Setting up new database '%v'\n", *dbFile)
		config.InitConfig()
	}
	if *passphrase != "" {
		config.SetPassphrase(*passphrase)
	}
	if *pageTitle != "" {
		config.SetPageTitle(*pageTitle)
	}
	defer database.Close()
	log.Printf("Server started listening on port %v\n", *port)
	smhttp.StartServer(*port)
}
