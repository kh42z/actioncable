package actioncable

import (
	"io/ioutil"
	"log"
	"sync"
)

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

func (ac *Client) Run() error {
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

func WithLogger(logger *log.Logger) Option {
	return func(c *Client) {
		c.logger = logger
	}
}

type Client struct {
	ws       JSONReadWriter
	channels map[string]ChannelHandler
	emit     chan *message
	stop     chan struct{}
	once     sync.Once
	logger   *log.Logger
}

type JSONReadWriter interface {
	ReadJSON(v interface{}) error
	WriteJSON(v interface{}) error
}

type Option func(*Client)
