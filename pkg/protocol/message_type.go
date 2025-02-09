package protocol

type MessageType string

const (
	ChallengeMessageType MessageType = "challenge"
	SolutionMessageType  MessageType = "solution"
	ErrorMessageType     MessageType = "error"
)

func (t MessageType) String() string {
	return string(t)
}
