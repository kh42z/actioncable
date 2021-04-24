package actioncable

import (
	"io/ioutil"
	"log"
	"sync"
)

func NewClient(ws JSONReadWriter, opts ...Option) *Client {
	c := &Client{
		ws:       ws,
		emit:     make(chan *Message),
		quit:     make(chan struct{}),
		channels: make(map[string]ChannelCallbacker),
		logger:   log.New(ioutil.Discard, "ActionCable: ", log.LstdFlags),
		onStartup: func() {
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
		ac.once.Do(func() { close(ac.quit) })
		return err
	}
	return nil
}

func WithLogger(logger *log.Logger) Option {
	return func(c *Client) {
		c.logger = logger
	}
}

func WithOnStartup(fn func()) Option {
	return func(c *Client) {
		c.onStartup = fn
	}
}

type Client struct {
	ws        JSONReadWriter
	emit      chan *Message
	quit      chan struct{}
	channels  map[string]ChannelCallbacker
	once      sync.Once
	logger    *log.Logger
	onStartup func()
}

type JSONReadWriter interface {
	ReadJSON(v interface{}) error
	WriteJSON(v interface{}) error
}

type Option func(*Client)
