package natsjs

import (
	"fmt"

	"github.com/nats-io/nats.go"
	"github.com/nats-io/nats.go/jetstream"
)

type Options struct {
	URL        string
	User       string
	Password   string
	ClientName string

	Verbose              bool
	AllowReconnect       bool
	RetryOnFailedConnect bool

	PublishAsyncMaxPending int
}

type Client struct {
	conn *nats.Conn
	js   jetstream.JetStream
}

func NewClient(opts Options) (*Client, error) {
	nc, err := nats.Connect(
		opts.URL,
		nats.Name(opts.ClientName),
		nats.UserInfo(opts.User, opts.Password),
		func(o *nats.Options) error {
			o.Verbose = opts.Verbose
			o.AllowReconnect = opts.AllowReconnect
			o.RetryOnFailedConnect = opts.RetryOnFailedConnect
			return nil
		},
	)
	if err != nil {
		return nil, fmt.Errorf("nats connect: %w", err)
	}

	js, err := jetstream.New(nc, jetstream.WithPublishAsyncMaxPending(opts.PublishAsyncMaxPending))
	if err != nil {
		nc.Close()
		return nil, fmt.Errorf("jetstream new: %w", err)
	}

	return &Client{conn: nc, js: js}, nil
}

func (c *Client) JS() jetstream.JetStream {
	return c.js
}

func (c *Client) Conn() *nats.Conn {
	return c.conn
}

func (c *Client) Close() {
	if c.conn != nil {
		c.conn.Close()
	}
}
