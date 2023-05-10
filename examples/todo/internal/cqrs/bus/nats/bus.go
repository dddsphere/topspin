package nats

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"io/ioutil"
	"time"

	"github.com/nats-io/nats.go"

	"github.com/dddsphere/topspin"
	"github.com/dddsphere/topspin/examples/todo/internal/cqrs"
)

const (
	commandSubject = "command"
	querySubject   = "query"
	eventSubject   = "event"
)

type (
	// BusManager is an implementation of bus.Manager built on top of NATS.
	BusManager struct {
		*topspin.SimpleWorker
		config *topspin.Config
		cqrs   *cqrs.Manager
		nats   *Client
	}
)

func NewBusManager(name string, cfg *topspin.Config, nats *Client, cqrs *cqrs.Manager, log topspin.Logger) *BusManager {
	return &BusManager{
		SimpleWorker: topspin.NewWorker(name, log),
		config:       cfg,
		cqrs:         cqrs,
		nats:         nats,
	}
}

func (bm *BusManager) Start() (err error) {
	err = bm.nats.Start()
	if err != nil {
		return fmt.Errorf("BusManager start error: %w", err)
	}

	bm.nats.Subscribe(commandSubject, bm.cqrs.HandleFunc())
	return nil
}

func (bm *BusManager) sampleHandler() nats.MsgHandler {
	return func(m *nats.Msg) {
		buf := bytes.NewBuffer(m.Data)
		dec := gob.NewDecoder(buf)

		ce := cqrs.CommandEvent{}
		err := dec.Decode(&ce)
		if err != nil {
			bm.Log().Errorf("Cannot decode command event: %s", err.Error())
		}

		bm.Log().Infof("Received a command event with ID: %s", ce.TracingID)
	}
}

func (bm *BusManager) SendCommand(cmd string, payload []byte, tracingID string) error {
	ce, err := encodeCommand(cmd, payload, tracingID)
	if err != nil {
		bm.Log().Errorf("Send command error: %s", err.Error())
		return err
	}

	return bm.nats.PublishEvent(commandSubject, ce)
}

func (bm *BusManager) SendEvent(event string, payload []byte, tracingID string) error {
	panic("not implemented")
}

func (bm *BusManager) Query(query string, payload []byte, timeout time.Duration, tracingID string) (response []byte, err error) {
	panic("not implemented")
}

// Helpers

func encodeCommand(cmd string, payload []byte, tracingID string) (cmdEvent []byte, err error) {
	ce := cqrs.CommandEvent{
		Command:   cmd,
		Payload:   payload,
		TracingID: tracingID,
	}

	buf := bytes.Buffer{}
	err = gob.NewEncoder(&buf).Encode(ce)
	if err != nil {
		return cmdEvent, err
	}

	cmdEvent, err = ioutil.ReadAll(&buf)
	if err != nil {
		return cmdEvent, err
	}

	return cmdEvent, err
}
