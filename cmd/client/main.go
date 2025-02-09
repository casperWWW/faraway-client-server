package main

import (
	"faraway/internal/client"
	"flag"
	"log"
)

func main() {
	serverAddr := flag.String("server", "localhost:8080", "Server address to connect to")
	flag.Parse()

	if err := runClient(*serverAddr); err != nil {
		log.Fatal(err)
	}
}

func runClient(serverAddr string) error {
	c := client.NewClient(serverAddr)
	return c.Connect()
}
