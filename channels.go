package actioncable

import (
	"encoding/json"
)

type ChannelHandler interface {
	SubscriptionHandler(*Client, int)
	MessageHandler(*Client, []byte, int)
}

func (ac *Client) AddChannelHandler(name string, event ChannelHandler) {
	ac.channels[name] = event
}

func (ac *Client) handleEvent(event *event) {
	var i identifier
	err := json.Unmarshal([]byte(event.Identifier), &i)
	if err != nil {
		return
	}
	for name, e := range ac.channels {
		if name == i.Channel {
			e.MessageHandler(ac, []byte(event.Message), i.ID)
		}
	}
}
