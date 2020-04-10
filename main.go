package main

import (
	"flag"
	"fmt"
	"log"
)

func main() {
	// TODO: Use getopt package instead of flags?
	port := flag.Int("p", 8080, "Port")
	dbFile := flag.String("d", "sermoni.db", "Database file")
	//password := flag.String("w", "", "Password for the web interface")
	flag.Parse()
	if err := initDB(*dbFile); err != nil {
		log.Fatal(err)
	}
	defer closeDB()
	fmt.Printf("Server running on port %v", *port)
}
