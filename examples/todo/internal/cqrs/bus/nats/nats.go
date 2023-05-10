package nats

import (
	"fmt"
	"runtime"

	"github.com/nats-io/nats.go"

	"github.com/dddsphere/topspin"
)

const (
	cfgHostKey = "nats.host"
	cfgPortKey = "nats.port"
)

type (
	Client struct {
		*topspin.SimpleWorker
		config *topspin.Config
		conn   *nats.Conn
	}
)

const (
	defaultHost = "localhost"
	defaultPort = 4222
)

func NewClient(name string, cfg *topspin.Config, log topspin.Logger) *Client {
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

	return nil
}

func (c *Client) address() (address string) {
	host := defaultHost
	if c.config.GetString(cfgHostKey) == "" {
		host = c.config.GetString(cfgHostKey)
	}

	port := defaultPort
	if c.config.GetInt(cfgPortKey) == 0 {
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

func (c *Client) Subscribe(subject string, handler nats.MsgHandler) {
	c.Log().Infof("NATS subscribed through: %s", c.conn.ConnectedAddr())

	var err error
	_, err = c.conn.Subscribe(subject, handler)
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
