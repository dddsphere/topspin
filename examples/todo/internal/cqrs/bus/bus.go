package bus

import "time"

type (
	Manager interface {
		SendCommand(cmd string, payload []byte, tracingID string) error
		SendEvent(event string, payload []byte, tracingID string) error
		Query(query string, payload []byte, timeout time.Duration, tracingID string) (response []byte, err error)
	}
)
