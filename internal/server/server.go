package server

import (
	"encoding/gob"
	"faraway/pkg/pow"
	"faraway/pkg/protocol"
	"faraway/pkg/quotes"
	"fmt"
	"log"
	"net"
	"time"
)

const (
	readTimeout  = 30 * time.Second
	writeTimeout = 30 * time.Second
)

type Server struct {
	addr string
}

type Message struct {
	Type    string
	Payload string
}

func NewServer(addr string) *Server {
	return &Server{addr: addr}
}

func (s *Server) Start() error {
	listener, err := net.Listen("tcp", s.addr)
	if err != nil {
		return fmt.Errorf("failed to start server: %v", err)
	}
	defer listener.Close()

	log.Printf("Server listening on %s", s.addr)

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Printf("Failed to accept connection: %v", err)
			continue
		}
		go s.handleConnection(conn)
	}
}

func (s *Server) handleConnection(conn net.Conn) {
	defer func(conn net.Conn) {
		err := conn.Close()
		if err != nil {
			log.Printf("Failed to close connection: %v", err)
		}
	}(conn)

	encoder := gob.NewEncoder(conn)
	decoder := gob.NewDecoder(conn)

	// Setup write deadline
	if err := conn.SetWriteDeadline(time.Now().Add(writeTimeout)); err != nil {
		log.Printf("Failed to set write deadline: %v", err)
		return
	}

	// Send challenge
	challenge := pow.NewChallenge("word-of-wisdom")
	if err := encoder.Encode(Message{Type: protocol.ChallengeMessageType.String(), Payload: challenge.Encode()}); err != nil {
		log.Printf("Failed to send challenge: %v", err)
		return
	}

	// Setup read deadline
	if err := conn.SetReadDeadline(time.Now().Add(readTimeout)); err != nil {
		log.Printf("Failed to set read deadline: %v", err)
		return
	}

	// Receive solution
	var msg Message
	if err := decoder.Decode(&msg); err != nil {
		log.Printf("Failed to receive solution: %v", err)
		return
	}

	if msg.Type != protocol.SolutionMessageType.String() {
		log.Printf("Unexpected message type: %s", msg.Type)
		return
	}

	// Verify solution
	solution := pow.Solution{
		Challenge: challenge,
		Counter:   uint64(0),
	}
	_, err := fmt.Sscanf(msg.Payload, "%d", &solution.Counter)
	if err != nil {
		log.Printf("Failed to parse solution: %v", err)
		return
	}

	if !solution.Verify() {
		err := encoder.Encode(Message{Type: protocol.ErrorMessageType.String(), Payload: "Invalid solution"})
		if err != nil {
			log.Printf("Failed to send solution: %v", err)
			return
		}
		return
	}

	// Send quote
	quote := quotes.GetRandomQuote()
	err = encoder.Encode(Message{Type: "quote", Payload: quote})
	if err != nil {
		log.Printf("Failed to send quote: %v", err)
		return
	}
}
