[![action-cable](https://github.com/kh42z/actioncable/actions/workflows/workflow.yml/badge.svg)](https://github.com/kh42z/actioncable/actions/workflows/workflow.yml)
[![Coverage Status](https://coveralls.io/repos/github/kh42z/actioncable/badge.svg?branch=master)](https://coveralls.io/github/kh42z/actioncable?branch=master)
[![Go Reference](https://pkg.go.dev/badge/github.com/kh42z/actioncable.svg)](https://pkg.go.dev/github.com/kh42z/actioncable)
# actioncable

A Go Client for Rails websocket ActionCable

## Usage

```go
package main

import (
	"github.com/kh42z/actioncable"
	"github.com/gorilla/websocket"
	"log"
)

func main() {
	ws, _, err := websocket.DefaultDialer.Dial("ws://localhost:3000/cable", nil)
	if err != nil {
		log.Fatalln("unable to connect ", err)
	}
	defer ws.Close()
	c := actioncable.NewClient(ws))
	c.AddChannelHandler("UserChannel", &UserChannel{})
	go func() {
		c.Subscribe("UserChannel", 1)
	}()
	if err := c.Run(); err != nil {
		log.Fatal("Actioncable: ", err)
	}
}

func (u *UserChannel) SubscriptionHandler(c *actioncable.Client, id int) {
	log.Println("UserChannel, successfully subscribed to id: ", id)
	c.SendMessage("UserChannel", id, "Hello!")
}

func (u *UserChannel) MessageHandler(_ *actioncable.Client, content []byte, id int) {
	log.Printf("UserChannel_%d: <%s>", id, content)
}

type UserChannel struct{}
```
