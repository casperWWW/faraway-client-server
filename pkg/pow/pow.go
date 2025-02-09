package pow

import (
	"crypto/sha256"
	"encoding/base64"
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"strings"
	"time"
)

const (
	// DefaultDifficulty is the number of leading zero bits required
	DefaultDifficulty = 20
	// DefaultValidity is how long the challenge is valid
	DefaultValidity = 2 * time.Minute
)

// Challenge represents a proof of work challenge
type Challenge struct {
	Timestamp  int64
	Resource   string
	Difficulty int
	Rand       []byte
}

// Solution represents a solution to a proof of work challenge
type Solution struct {
	Challenge Challenge
	Counter   uint64
}

// NewChallenge creates a new proof of work challenge
func NewChallenge(resource string) Challenge {
	randBytes := make([]byte, 8)
	binary.LittleEndian.PutUint64(randBytes, uint64(time.Now().UnixNano()))

	return Challenge{
		Timestamp:  time.Now().Unix(),
		Resource:   resource,
		Difficulty: DefaultDifficulty,
		Rand:       randBytes,
	}
}

// Encode returns the base64 encoded string representation of the challenge
func (c Challenge) Encode() string {
	return base64.StdEncoding.EncodeToString([]byte(fmt.Sprintf("%d:%s:%d:%x",
		c.Timestamp, c.Resource, c.Difficulty, c.Rand)))
}

// DecodeChallenge decodes a base64 encoded challenge string
func DecodeChallenge(s string) (Challenge, error) {
	data, err := base64.StdEncoding.DecodeString(s)
	if err != nil {
		return Challenge{}, err
	}

	parts := strings.Split(string(data), ":")
	if len(parts) != 4 {
		return Challenge{}, fmt.Errorf("invalid challenge format")
	}

	var timestamp int64
	if _, err := fmt.Sscanf(parts[0], "%d", &timestamp); err != nil {
		return Challenge{}, fmt.Errorf("invalid timestamp: %v", err)
	}

	var difficulty int
	if _, err := fmt.Sscanf(parts[2], "%d", &difficulty); err != nil {
		return Challenge{}, fmt.Errorf("invalid difficulty: %v", err)
	}

	randBytes, err := hex.DecodeString(parts[3])
	if err != nil {
		return Challenge{}, fmt.Errorf("invalid random bytes: %v", err)
	}

	return Challenge{
		Timestamp:  timestamp,
		Resource:   parts[1],
		Difficulty: difficulty,
		Rand:       randBytes,
	}, nil
}

// Verify checks if the solution is valid
func (s Solution) Verify() bool {
	// Check if the challenge has expired
	if time.Now().Unix()-s.Challenge.Timestamp > int64(DefaultValidity.Seconds()) {
		return false
	}

	hash := s.computeHash()
	return countLeadingZeros(hash) >= s.Challenge.Difficulty
}

// Solve finds a solution to the challenge
func (c Challenge) Solve() Solution {
	s := Solution{Challenge: c}
	for {
		if countLeadingZeros(s.computeHash()) >= c.Difficulty {
			return s
		}
		s.Counter++
	}
}

func (s Solution) computeHash() []byte {
	h := sha256.New()
	h.Write([]byte(fmt.Sprintf("%d:%s:%d:%x:%d",
		s.Challenge.Timestamp,
		s.Challenge.Resource,
		s.Challenge.Difficulty,
		s.Challenge.Rand,
		s.Counter)))
	return h.Sum(nil)
}

func countLeadingZeros(data []byte) int {
	totalBits := len(data) * 8
	for i := 0; i < totalBits; i++ {
		if data[i/8]&(1<<uint(7-i%8)) != 0 {
			return i
		}
	}
	return totalBits
}
