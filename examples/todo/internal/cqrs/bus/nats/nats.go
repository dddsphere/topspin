package nats

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"runtime"

	"github.com/nats-io/nats.go"

	"github.com/dddsphere/topspin"
	"github.com/dddsphere/topspin/examples/todo/internal/config"
)

type (
	Client struct {
		*topspin.SimpleWorker
		config *config.Config
		conn   *nats.Conn
	}
)

const (
	defaultHost = "localhost"
	defaultPort = 4222
)

func NewClient(name string, cfg *config.Config, log topspin.Logger) *Client {
	return &Client{
		SimpleWorker: topspin.NewWorker(name, log),
		config:       cfg,
	}
}

func (c *Client) Start() error {
	c.Log().Infof("NATS client connecting to %s", c.address())

	var err error
	c.conn, err = nats.Connect(c.address())
	if err != nil {
		return fmt.Errorf("nats connection cannot be established: %w", err)
	}

	// Subscriptions
	// WIP: Not a client responsibility,
	// Move this up, just only to verify subscriptions are working.
	c.Subscribe(commandSubject)

	return nil
}

func (c *Client) address() (address string) {
	host := defaultHost
	if c.config.NATS.Host == "" {
		host = c.config.NATS.Host
	}

	port := defaultPort
	if c.config.NATS.Port == 0 {
		port = defaultPort
	}

	return fmt.Sprintf("nats://%s:%d", host, port)
}

func (c *Client) PublishEvent(subject string, commandEvent []byte) error {
	c.Log().Infof("NATS publishing through: %s", c.conn.ConnectedAddr())

	err := c.conn.Publish(subject, commandEvent)
	if err != nil {
		return fmt.Errorf("NATS client error: %w", err)
	}
	return nil
}

func (c *Client) Subscribe(subject string) {
	c.Log().Infof("NATS subscribed through: %s", c.conn.ConnectedAddr())

	var err error
	_, err = c.conn.Subscribe(subject, func(m *nats.Msg) {
		buf := bytes.NewBuffer(m.Data)
		dec := gob.NewDecoder(buf)

		ce := CommandEvent{}
		err := dec.Decode(&ce)
		if err != nil {
			c.Log().Errorf("Cannot decode command event: %s", err.Error())
		}

		c.Log().Infof("Received a command event with ID: %s", ce.TracingID)
	})

	if err != nil {
		c.Log().Errorf("NATS command subscription error: %s", err.Error())
	}

	err = c.conn.Flush()
	if err != nil {
		c.Log().Errorf("NATS flush error: %s", err.Error())
	}

	err = c.conn.LastError()
	if err != nil {
		c.Log().Error(err.Error())
	}

	c.Log().Infof("Listening on '%s' subject", subject)

	runtime.Goexit()
}
