package cqrs

import (
	"bytes"
	"encoding/gob"
	"encoding/json"

	"github.com/nats-io/nats.go"

	"github.com/dddsphere/topspin"
	"github.com/dddsphere/topspin/examples/todo/internal/cqrs/command"
)

type (
	Manager struct {
		topspin.Worker
		*topspin.CQRSManager
	}
)

type (
	CommandEvent struct {
		Command   string
		Payload   []byte
		TracingID string
	}
)

func NewManager(log topspin.Logger) *Manager {
	return &Manager{
		Worker:      topspin.NewWorker("CQRS-Manager", log),
		CQRSManager: topspin.NewCQRSManager(),
	}
}

func (cm *Manager) HandleFunc() nats.MsgHandler {
	return cm.handle
}

// WIP: Lots to do here yet
func (cm *Manager) handle(msg *nats.Msg) {
	ce := CommandEvent{}
	err := gob.NewDecoder(bytes.NewReader(msg.Data)).Decode(&ce)
	if err != nil {
		cm.Log().Errorf("Cannot decode command event: %w", err)
	}

	cm.Log().Infof("Received a command event with ID: %s", ce.TracingID)

	var data command.CreateListCommandData

	err = json.Unmarshal(ce.Payload, &data)
	if err != nil {
		cm.Log().Errorf("Cannot unmarshall command event payload: %w", err)
	}

	cm.Log().Debugf("Processing '%s' command with this data: %+v", ce.Command, data)

	//err = cm.todoService.CreateList(context.Background(), data.Name, data.Description)
	if err != nil {
		//return fmt.Errorf("%s handle error: %w", cmd.Name(), err)
		cm.Log().Errorf("Cannot process command '%s': %w", ce.Command, err)
	}
}
