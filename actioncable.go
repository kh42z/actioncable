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
		quit:     make(chan struct{}),
		channels: make(map[string]ChannelCallbacker),
		logger:   log.New(ioutil.Discard, "actionCable: ", log.LstdFlags),
		onWelcome: func() {
		},
	}
	for _, opt := range opts {
		opt(c)
	}
	return c
}

func (ac *Client) Start() error {
	go ac.send()
	if err := ac.receive(); err != nil {
		ac.stop()
		return err
	}
	return nil
}

func (ac *Client) stop() {
	ac.once.Do(func() { close(ac.quit) })
}

func WithLogger(logger *log.Logger) Option {
	return func(c *Client) {
		c.logger = logger
	}
}

func WithOnWelcome(fn func()) Option {
	return func(c *Client) {
		c.onWelcome = fn
	}
}

type Client struct {
	ws        JSONReadWriter
	emit      chan *message
	quit      chan struct{}
	channels  map[string]ChannelCallbacker
	once      sync.Once
	logger    *log.Logger
	onWelcome func()
}

type JSONReadWriter interface {
	ReadJSON(v interface{}) error
	WriteJSON(v interface{}) error
}

type Option func(*Client)
