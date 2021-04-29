package actioncable

import (
	"io/ioutil"
	"log"
	"sync"
)

// NewClient returns a *Client
func NewClient(ws JSONReadWriter, opts ...Option) *Client {
	c := &Client{
		ws:       ws,
		emit:     make(chan *message),
		channels: make(map[string]ChannelHandler),
		stop:     make(chan struct{}),
		logger:   log.New(ioutil.Discard, "actionCable: ", log.LstdFlags),
	}
	for _, opt := range opts {
		opt(c)
	}
	return c
}

// Run starts reading and writing on the Connection provided
func (ac *Client) Run() error {
	if err := ac.waitWelcome(); err != nil {
		return err
	}
	go ac.send()
	if err := ac.receive(); err != nil {
		ac.exit()
		return err
	}
	return nil
}

func (ac *Client) exit() {
	ac.once.Do(func() { close(ac.stop) })
}

// WithLogger returns an Option used by NewClient to set Client logger
func WithLogger(logger *log.Logger) Option {
	return func(c *Client) {
		c.logger = logger
	}
}

// Client type represents an Actioncable Client
type Client struct {
	ws       JSONReadWriter
	channels map[string]ChannelHandler
	emit     chan *message
	stop     chan struct{}
	once     sync.Once
	logger   *log.Logger
}

// JSONReadWriter represents a websocket implementation. github.com/gorilla/websocket *Conn satisfies this interface.
type JSONReadWriter interface {
	ReadJSON(v interface{}) error
	WriteJSON(v interface{}) error
}

// Option is a func that modifies Client values
type Option func(*Client)
