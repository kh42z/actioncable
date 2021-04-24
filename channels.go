package actioncable

import (
	"encoding/json"
)

type ChannelCallbacker interface {
	OnSubscription(int)
	OnMessage([]byte, int)
}

func (ac *Client) RegisterChannelCallbacks(name string, event ChannelCallbacker) {
	ac.channels[name] = event
}

func (ac *Client) handleEvent(event *event) {
	var i identifier
	err := json.Unmarshal([]byte(event.Identifier), &i)
	if err != nil {
		ac.logger.Println("handleEvent : unable to unmarshal Identifier", i)
		return
	}
	for name, e := range ac.channels {
		if name == i.Channel {
			e.OnMessage(event.Message, i.ID)
		}
	}
}
