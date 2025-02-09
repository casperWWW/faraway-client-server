package main

import (
	"faraway/internal/server"
	"flag"
	"log"
)

func main() {
	addr := flag.String("addr", ":8080", "Server address to listen on")
	flag.Parse()

	if err := runServer(*addr); err != nil {
		log.Fatal(err)
	}
}

func runServer(addr string) error {
	s := server.NewServer(addr)
	return s.Start()
}
