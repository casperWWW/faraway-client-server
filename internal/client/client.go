package client

import (
	"encoding/gob"
	"faraway/pkg/pow"
	"faraway/pkg/protocol"
	"fmt"
	"net"
	"time"
)

const (
	readTimeout  = 30 * time.Second
	writeTimeout = 30 * time.Second
)

type Client struct {
	serverAddr string
}

type Message struct {
	Type    string
	Payload string
}

func NewClient(serverAddr string) *Client {
	return &Client{serverAddr: serverAddr}
}

func (c *Client) Connect() error {
	conn, err := net.Dial("tcp", c.serverAddr)
	if err != nil {
		return fmt.Errorf("failed to connect to server: %v", err)
	}
	defer conn.Close()

	encoder := gob.NewEncoder(conn)
	decoder := gob.NewDecoder(conn)

	// Receive challenge
	var msg Message
	if err := decoder.Decode(&msg); err != nil {
		return fmt.Errorf("failed to receive challenge: %v", err)
	}

	if msg.Type != string(protocol.ChallengeMessageType) {
		return fmt.Errorf("unexpected message type: %s", msg.Type)
	}

	// Decode and solve challenge
	challenge, err := pow.DecodeChallenge(msg.Payload)
	if err != nil {
		return fmt.Errorf("failed to decode challenge: %v", err)
	}

	solution := challenge.Solve()

	// Send solution
	if err := encoder.Encode(Message{Type: string(protocol.SolutionMessageType), Payload: fmt.Sprintf("%d", solution.Counter)}); err != nil {
		return fmt.Errorf("failed to send solution: %v", err)
	}

	// Receive quote
	if err := decoder.Decode(&msg); err != nil {
		return fmt.Errorf("failed to receive quote: %v", err)
	}

	if msg.Type == "error" {
		return fmt.Errorf("server error: %s", msg.Payload)
	}

	if msg.Type != "quote" {
		return fmt.Errorf("unexpected message type: %s", msg.Type)
	}

	fmt.Println("Received quote:", msg.Payload)
	return nil
}
