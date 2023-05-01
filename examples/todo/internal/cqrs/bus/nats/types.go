package nats

type (
	CommandEvent struct {
		Command   string
		Payload   []byte
		TracingID string
	}
)
